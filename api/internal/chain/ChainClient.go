package chain

import (
	"context"
	"math"
	"math/big"
	"strings"
	"sync"
	"time"
	wowfish "wowfish/api/contract"
	"wowfish/api/internal/callback"
	"wowfish/api/internal/config"
	"wowfish/api/pkg/kms"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zeromicro/go-zero/core/logx"
)

type Tx struct {
	To    string
	Value uint64
	Data  []byte
}

type ChainClient struct {
	Provider   *ethclient.Client
	Kms        *kms.KmsEthMgr
	chainId    *big.Int
	ctx        context.Context
	rpcAddress string
}

var once sync.Once
var instance = &ChainClient{}

func ChainClientInstance() *ChainClient {
	once.Do(func() {
		instance = &ChainClient{}
	})
	return instance
}

func (this *ChainClient) connectWsServer(ctx context.Context) {
	provider, err := ethclient.DialContext(ctx, this.rpcAddress)
	for {
		if err == nil {
			break
		}
		logx.Errorf("connect rpc %s error %s, reconnect....", this.rpcAddress, err.Error())
		time.Sleep(time.Duration(1) * time.Second)
		provider, err = ethclient.DialContext(ctx, this.rpcAddress)
	}
	logx.Infof("connect to ws %s sucess", this.rpcAddress)
	this.Provider = provider
}

func (this *ChainClient) InitChainClient(ctx context.Context, config *config.Config) error {
	this.rpcAddress = config.Chain.Rpc
	this.connectWsServer(ctx)
	chainId, err := this.Provider.ChainID(ctx)
	if err != nil {
		logx.Errorf("get chain id error %s", err.Error())
		return err
	}

	this.chainId = chainId
	this.ctx = ctx
	this.Kms = kms.NewKmsInstance()
	alias := []string{config.Kms.GameWallet, config.Kms.GantWallet, config.Kms.DisposeWallet}
	err = this.Kms.InitKmsEthMgr(config.Kms.Key, config.Kms.Secret, alias)
	if err != nil {
		logx.Errorf("int kms manager error %s", err.Error())
		return err
	}
	tokens := []string{config.Chain.WowToken}
	this.watchToken(tokens)
	return nil
}

func (this *ChainClient) Exit() {
	this.Provider.Close()
}

func (this *ChainClient) CommitTranscation(tx *Tx, from string) (string, error) {
	//make raw tx
	nonce, err := this.Provider.PendingNonceAt(this.ctx, common.HexToAddress(from))
	if err != nil {
		logx.Errorf("pending nonce error %s", err.Error())
		return "", err
	}

	value := big.NewInt(int64(tx.Value)) // in wei (1 eth)
	gasPrice, err := this.Provider.SuggestGasPrice(this.ctx)
	if err != nil {
		logx.Errorf("suggest gas price error %s", err.Error())
		return "", err
	}

	toAddress := common.HexToAddress(tx.To)

	eGas, err := this.GetEstimateGas(tx, from)
	if err != nil {
		logx.Errorf("estimate gas error %s", err.Error())
		return "", err
	}

	dynamicTx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   this.chainId,
		Nonce:     nonce,
		GasFeeCap: gasPrice,
		GasTipCap: big.NewInt(0),
		Gas:       eGas,
		To:        &toAddress,
		Value:     value,
		Data:      tx.Data,
	})

	signer := types.LatestSignerForChainID(this.chainId)
	txHashBytes := signer.Hash(dynamicTx).Bytes()

	txSignResult, err := this.Kms.Sign(from, txHashBytes, "DIGEST")
	if err != nil {
		logx.Errorf("sign dynamic tx error %s", err.Error())
		return "", err
	}
	dynamicTx, err = dynamicTx.WithSignature(signer, txSignResult)
	if err != nil {
		logx.Errorf("tx WithSignature error %s", err.Error())
		return "", err
	}
	err = this.Provider.SendTransaction(this.ctx, dynamicTx)
	if err != nil {
		logx.Errorf("send transcation tx error %s", err.Error())
		return "", err
	}
	return dynamicTx.Hash().Hex(), nil
}

func (this *ChainClient) GenTranction(from string) (*bind.TransactOpts, error) {
	txOp, err := this.Kms.GetNewKmsTransactor(this.ctx, this.chainId, from)
	if err != nil {
		logx.Errorf("gen kms transcatior error %s", err.Error())
		return nil, err
	}

	nonce, err := this.Provider.PendingNonceAt(this.ctx, txOp.From)
	if err != nil {
		logx.Errorf("pending nonce error %s", err.Error())
		return nil, err
	}

	gasPrice, err := this.Provider.SuggestGasPrice(this.ctx)
	if err != nil {
		logx.Errorf("suggest gas price error %s", err.Error())
		return nil, err
	}
	txOp.Nonce = big.NewInt(int64(nonce))
	txOp.GasFeeCap = gasPrice
	txOp.GasTipCap = big.NewInt(0)
	txOp.Signer = func(address common.Address, tx *types.Transaction) (*types.Transaction, error) {

		if address != txOp.From {
			return nil, bind.ErrNotAuthorized
		}
		signer := types.LatestSignerForChainID(this.chainId)
		txHashBytes := signer.Hash(tx).Bytes()
		signature, err := this.Kms.Sign(from, txHashBytes, "DIGEST")
		if err != nil {
			logx.Errorf("sign dynamic tx error %s", err.Error())
			return nil, err
		}
		return tx.WithSignature(signer, signature)
	}
	return txOp, nil
}

func (this *ChainClient) GetEstimateGas(tx *Tx, address string) (eGas uint64, err error) {
	toAddress := common.HexToAddress(tx.To)
	cMsg := ethereum.CallMsg{
		From: common.HexToAddress(address),
		To:   &toAddress,
		Data: tx.Data,
	}

	eGas, err = this.Provider.EstimateGas(this.ctx, cMsg)
	if err != nil {
		logx.Errorf("EstimateGas err:%s", err.Error())
		return 0, err
	}
	return eGas, nil
}

func (this *ChainClient) watchToken(tokens []string) {
	query := ethereum.FilterQuery{
		Addresses: []common.Address{},
		Topics:    [][]common.Hash{},
	}
	logTransferSig := []byte("Transfer(address,address,uint256)")
	logTransferSigHash := crypto.Keccak256Hash(logTransferSig)
	logx.Infof("subscript topic: %s", logTransferSigHash.Hex())
	decimals := make(map[string]uint8)
	for _, token := range tokens {

		tokenAddress := common.HexToAddress(token)
		//query decimal
		wowToken, err := wowfish.NewWowFishToken(tokenAddress, this.Provider)
		if err != nil {
			logx.Errorf("SubscribeFilterLogs create wowfish token error %s", err.Error())
			break
		}
		decimal, err := wowToken.Decimals(&bind.CallOpts{})
		if err != nil {
			logx.Errorf("SubscribeFilterLogs query Decimals error %s", err.Error())
			break
		}
		decimals[strings.ToLower(token)] = decimal
		logx.Infof("subscript address: %s", token)
		query.Addresses = append(query.Addresses, tokenAddress)
		query.Topics = append(query.Topics, []common.Hash{logTransferSigHash})
	}

	logs := make(chan types.Log)
	sub, err := this.Provider.SubscribeFilterLogs(this.ctx, query, logs)
	if err != nil {
		logx.Errorf("SubscribeFilterLogs error %s", err.Error())
		return
	}
	logx.Infof("subscribe query success %+v", query)
	contractAbi, err := wowfish.WowFishTokenMetaData.GetAbi()
	if err != nil {
		logx.Errorf("SubscribeFilterLogs error %s", err.Error())
		return
	}
	go func(d map[string]uint8) {
		for {
			select {
			case err := <-sub.Err():
				logx.Errorf("SubscribeFilterLogs error %s", err.Error())
				//need reconnect
				this.connectWsServer(this.ctx)
				//re watch
				sub, err = this.Provider.SubscribeFilterLogs(this.ctx, query, logs)
				if err != nil {
					logx.Errorf("SubscribeFilterLogs error %s", err.Error())
					return
				}
				logx.Infof("re subscribe query success %+v", query)
			case vLog := <-logs:
				transferContent, err := contractAbi.Unpack("Transfer", vLog.Data)
				if err != nil {
					logx.Error(err)
				}

				decimal := d[strings.ToLower(vLog.Address.Hex())]

				if decimal == 0 {
					logx.Errorf("revice event but not find decimal %s", vLog.Address.Hex())
					continue
				}

				content, ok := transferContent[0].(*big.Int)
				if !ok {
					logx.Error("read transfer content error")
					continue
				}
				x := math.Pow10(int(decimal))

				contentValue := big.NewFloat(float64(content.Uint64()))
				contentValue.Quo(contentValue, big.NewFloat(x))

				from := common.HexToAddress(vLog.Topics[1].Hex()).String()
				to := common.HexToAddress(vLog.Topics[2].Hex()).String()
				amount := contentValue.String()
				logx.Infof("From:%s,To:%s,Amount:%s",
					from,
					to,
					amount)
				//callback
				callback.Instance().Callback(&callback.CallBackToWowfishData{
					From:   from,
					To:     to,
					Amount: amount,
					Ret:    0,
				})
			}
		}
	}(decimals)

}
