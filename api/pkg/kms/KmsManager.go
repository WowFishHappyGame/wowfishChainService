package kms

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/asn1"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"math/big"
	"strings"
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

type KmsAlias struct {
	Kid  string
	Type string
}

type KMSWallet struct {
	Id      KmsAlias
	pubKey  *ecdsa.PublicKey //key public key
	Address common.Address
}

type KmsEthMgr struct {
	awsCfg           aws.Config
	kmsSvc           *kms.Client
	Wallets          map[string]KMSWallet //key id
	SigningAlgorithm string               // aws.String("ECDSA_SHA_256"),
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

func (m *KmsEthMgr) InitKmsEthMgr(kmsKey string, kmsSecret string, keyids []KmsAlias) (err error) {
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
	m.Wallets = make(map[string]KMSWallet)

	for _, id := range keyids {
		key := KMSWallet{}
		key.Id = id
		key.pubKey, err = ethawskmssigner.GetPubKey(m.kmsSvc, id.Kid)
		if err != nil {
			logx.Error("GetPubKey fail ", err.Error())
			return err
		}
		key.Address = crypto.PubkeyToAddress(*key.pubKey)
		addressStr := key.Address.Hex()
		m.Wallets[strings.ToLower(addressStr)] = key

		logx.Infof("Init KmsEthMgr, addr= %s, alias=%s", addressStr, id.Kid)
	}
	return nil
}

func (m *KmsEthMgr) GetNewKmsTransactor(ctx context.Context, chainID *big.Int, address string) (*bind.TransactOpts, error) {
	wallet, err := m.GetWalletInfo(address)
	if err != nil {
		return nil, err
	}
	op, err := ethawskmssigner.NewAwsKmsTransactorWithChainIDCtx(ctx, m.kmsSvc, wallet.Id.Kid, chainID)
	op.From = wallet.Address
	return op, err
}

func (m *KmsEthMgr) GetWalletInfo(address string) (*KMSWallet, error) {
	wallet, ok := m.Wallets[strings.ToLower(address)]
	if !ok {
		//logx.Error("waller address error in sign", address)
		return nil, errors.New("address not in kms")
	}
	return &wallet, nil
}

func (m *KmsEthMgr) Sign(address string, message []byte, msgType string) ([]byte, error) {

	wallet, err := m.GetWalletInfo(address)
	if err != nil {
		return nil, err
	}

	input := &kms.SignInput{
		KeyId:            aws.String(wallet.Id.Kid),  //aws.String(m.keyID), // aws.String("alias/aa-relayer-key"),
		Message:          message,                    // []byte("<message to be signed>")
		MessageType:      types.MessageType(msgType), // aws.String("RAW"), //aws.String("DIGEST"),
		SigningAlgorithm: "ECDSA_SHA_256",            // aws.String("ECDSA_SHA_384"),
	}

	signOut, err := m.kmsSvc.Sign(context.Background(), input)
	if err != nil {
		logx.Error("Sign err:", err.Error())
		return nil, fmt.Errorf("sign err:%s", err.Error())
	}

	var sigAsn1 asn1EcSig
	_, err = asn1.Unmarshal(signOut.Signature, &sigAsn1)
	if err != nil {
		return nil, fmt.Errorf("sign unmarshal err:%s", err.Error())
	}

	rBytes, sBytes := sigAsn1.R.Bytes, sigAsn1.S.Bytes

	// Adjust S value from signature according to Ethereum standard
	sBigInt := new(big.Int).SetBytes(sBytes)
	if sBigInt.Cmp(secp256k1HalfN) > 0 {
		sBytes = new(big.Int).Sub(secp256k1N, sBigInt).Bytes()
	}

	pubKeyBytes := secp256k1.S256().Marshal(wallet.pubKey.X, wallet.pubKey.Y)
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

func (m *KmsEthMgr) GetWalletAddress(t string) string {
	for k, v := range m.Wallets {
		if v.Id.Type == t {
			return k
		}
	}
	return ""
}
