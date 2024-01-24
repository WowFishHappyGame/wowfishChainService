package kms

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/asn1"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	ethawskmssigner "github.com/welthee/go-ethereum-aws-kms-tx-signer/v2"
	"github.com/zeromicro/go-zero/core/logx"
)

type asn1EcSig struct {
	R asn1.RawValue
	S asn1.RawValue
}

type KmsEthMgr struct {
	awsCfg           aws.Config
	kmsSvc           *kms.Client
	ctx              context.Context
	keyID            string           //key id
	Address          common.Address   //key wallet
	SigningAlgorithm string           // aws.String("ECDSA_SHA_256"),
	pubKey           *ecdsa.PublicKey //key public key
}

var secp256k1N = crypto.S256().Params().N
var secp256k1HalfN = new(big.Int).Div(secp256k1N, big.NewInt(2))

var kmsInstance *KmsEthMgr
var once sync.Once

func NewKmsInstance() *KmsEthMgr {
	once.Do(func() {
		kmsInstance = &KmsEthMgr{}
	})
	return kmsInstance
}

func (m *KmsEthMgr) InitKmsEthMgr(kmsKey string, kmsSecret string, keyid string) (err error) {
	m.awsCfg, err = config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     kmsKey,
				SecretAccessKey: kmsSecret,
			},
		}))
	if err != nil {
		log.Fatal(err)
	}
	m.awsCfg.Region = "ap-southeast-1"
	m.kmsSvc = kms.NewFromConfig(m.awsCfg)
	m.keyID = keyid

	m.pubKey, err = ethawskmssigner.GetPubKey(m.kmsSvc, m.keyID)
	if err != nil {
		logx.Error("GetPubKey fail ", err.Error())
		return err
	}

	m.Address = crypto.PubkeyToAddress(*m.pubKey)
	m.SigningAlgorithm = "ECDSA_SHA_256"

	logx.Info("Init KmsEthMgr, addr=", m.Address.Hex())

	return nil
}

func (m *KmsEthMgr) GetNewKmsTransactor(ctx context.Context, chainID *big.Int) (*bind.TransactOpts, error) {
	return ethawskmssigner.NewAwsKmsTransactorWithChainIDCtx(ctx, m.kmsSvc, m.keyID, chainID)
}

func (m *KmsEthMgr) Sign(message []byte, msgType string) ([]byte, error) {

	input := &kms.SignInput{
		KeyId:            aws.String(m.keyID),        //aws.String(m.keyID), // aws.String("alias/aa-relayer-key"),
		Message:          message,                    // []byte("<message to be signed>")
		MessageType:      types.MessageType(msgType), // aws.String("RAW"), //aws.String("DIGEST"),
		SigningAlgorithm: "ECDSA_SHA_256",            // aws.String("ECDSA_SHA_384"),
	}

	signOut, err := m.kmsSvc.Sign(context.Background(), input)
	if err != nil {
		logx.Error("Sign err:", err.Error())
		return nil, fmt.Errorf("Sign err:%s\n", err.Error())
	}

	var sigAsn1 asn1EcSig
	_, err = asn1.Unmarshal(signOut.Signature, &sigAsn1)
	if err != nil {
		return nil, fmt.Errorf("Sign unmarshal err:%s\n", err.Error())
	}

	rBytes, sBytes := sigAsn1.R.Bytes, sigAsn1.S.Bytes

	// Adjust S value from signature according to Ethereum standard
	sBigInt := new(big.Int).SetBytes(sBytes)
	if sBigInt.Cmp(secp256k1HalfN) > 0 {
		sBytes = new(big.Int).Sub(secp256k1N, sBigInt).Bytes()
	}

	pubKeyBytes := secp256k1.S256().Marshal(m.pubKey.X, m.pubKey.Y)
	signature, err := m.getEthereumSignature(pubKeyBytes, message, rBytes, sBytes)
	if err != nil {
		return nil, err
	}
	//fmt.Println(hexutil.Encode(signature))

	return signature, nil
}

func adjustSignatureLength(buffer []byte) []byte {
	buffer = bytes.TrimLeft(buffer, "\x00")
	for len(buffer) < 32 {
		zeroBuf := []byte{0}
		buffer = append(zeroBuf, buffer...)
	}
	return buffer
}

func (m *KmsEthMgr) getEthereumSignature(expectedPublicKeyBytes []byte, txHash []byte, r []byte, s []byte) ([]byte, error) {
	rsSignature := append(adjustSignatureLength(r), adjustSignatureLength(s)...)
	signature := append(rsSignature, []byte{0}...)

	recoveredPublicKeyBytes, err := crypto.Ecrecover(txHash, signature)
	if err != nil {
		return nil, err
	}

	if hex.EncodeToString(recoveredPublicKeyBytes) != hex.EncodeToString(expectedPublicKeyBytes) {
		signature = append(rsSignature, []byte{1}...)
		recoveredPublicKeyBytes, err = crypto.Ecrecover(txHash, signature)
		if err != nil {
			return nil, err
		}

		if hex.EncodeToString(recoveredPublicKeyBytes) != hex.EncodeToString(expectedPublicKeyBytes) {
			return nil, fmt.Errorf("can not reconstruct public key from sig")
		}
	}

	return signature, nil
}
