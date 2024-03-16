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

	et "wowfish/api/pkg/types"
)

type Tx struct {
	To    string
	Value uint64
	Data  []byte
}

type ChainClient struct {
	Provider    *ethclient.Client
	Kms         *kms.KmsEthMgr
	chainId     *big.Int
	ctx         context.Context
	rpcAddress  string
	bankAddress string
}

var once sync.Once
var instance = &ChainClient{}

func ChainClientInstance() *ChainClient {
	once.Do(func() {
		instance = &ChainClient{}
	})
	return instance
}

func (c *ChainClient) connectWsServer(ctx context.Context) {
	provider, err := ethclient.DialContext(ctx, c.rpcAddress)
	for {
		if err == nil {
			break
		}
		logx.Errorf("connect rpc %s error %s, reconnect....", c.rpcAddress, err.Error())
		time.Sleep(time.Duration(1) * time.Second)
		provider, err = ethclient.DialContext(ctx, c.rpcAddress)
	}
	logx.Infof("connect to ws %s sucess", c.rpcAddress)
	c.Provider = provider
}

func (c *ChainClient) InitChainClient(ctx context.Context, config *config.Config) error {
	c.rpcAddress = config.Chain.Rpc
	c.connectWsServer(ctx)
	chainId, err := c.Provider.ChainID(ctx)
	if err != nil {
		logx.Errorf("get chain id error %s", err.Error())
		return err
	}
	c.bankAddress = config.Chain.GameBank
	c.chainId = chainId
	c.ctx = ctx
	c.Kms = kms.NewKmsInstance()
	alias := []kms.KmsAlias{
		{
			Kid:  config.Kms.GameWallet,
			Type: et.GameWallet,
		},
		{
			Kid:  config.Kms.GainWallet,
			Type: et.GainWallet,
		},
		{
			Kid:  config.Kms.DisposeWallet,
			Type: et.DisposeWallet,
		},
		{
			Kid:  config.Kms.AdminWallet,
			Type: et.AdminWallet,
		},
	}

	err = c.Kms.InitKmsEthMgr(config.Kms.Key, config.Kms.Secret, alias)
	if err != nil {
		logx.Errorf("int kms manager error %s", err.Error())
		return err
	}
	tokens := []string{config.Chain.WowToken}

	var wallets []string
	wallets = append(wallets, kms.NewKmsInstance().GetWalletAddress(et.GameWallet))
	wallets = append(wallets, c.bankAddress)
	wallets = append(wallets, kms.NewKmsInstance().GetWalletAddress(et.DisposeWallet))
	wallets = append(wallets, kms.NewKmsInstance().GetWalletAddress(et.GainWallet))
	c.watchToken(tokens, wallets)
	return nil
}

func (c *ChainClient) Exit() {
	c.Provider.Close()
}

func (c *ChainClient) CommitTranscation(tx *Tx, from string) (string, error) {
	//make raw tx
	nonce, err := c.Provider.PendingNonceAt(c.ctx, common.HexToAddress(from))
	if err != nil {
		logx.Errorf("pending nonce error %s", err.Error())
		return "", err
	}

	value := big.NewInt(int64(tx.Value)) // in wei (1 eth)
	gasPrice, err := c.Provider.SuggestGasPrice(c.ctx)
	if err != nil {
		logx.Errorf("suggest gas price error %s", err.Error())
		return "", err
	}

	toAddress := common.HexToAddress(tx.To)

	eGas, err := c.GetEstimateGas(tx, from)
	if err != nil {
		logx.Errorf("estimate gas error %s", err.Error())
		return "", err
	}

	dynamicTx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   c.chainId,
		Nonce:     nonce,
		GasFeeCap: gasPrice,
		GasTipCap: big.NewInt(0),
		Gas:       eGas,
		To:        &toAddress,
		Value:     value,
		Data:      tx.Data,
	})

	signer := types.LatestSignerForChainID(c.chainId)
	txHashBytes := signer.Hash(dynamicTx).Bytes()

	txSignResult, err := c.Kms.Sign(from, txHashBytes, "DIGEST")
	if err != nil {
		logx.Errorf("sign dynamic tx error %s", err.Error())
		return "", err
	}
	dynamicTx, err = dynamicTx.WithSignature(signer, txSignResult)
	if err != nil {
		logx.Errorf("tx WithSignature error %s", err.Error())
		return "", err
	}
	err = c.Provider.SendTransaction(c.ctx, dynamicTx)
	if err != nil {
		logx.Errorf("send transcation tx error %s", err.Error())
		return "", err
	}
	return dynamicTx.Hash().Hex(), nil
}

func (c *ChainClient) GenTranction(from string) (*bind.TransactOpts, error) {
	txOp, err := c.Kms.GetNewKmsTransactor(c.ctx, c.chainId, from)
	if err != nil {
		logx.Errorf("gen kms transcatior error %s", err.Error())
		return nil, err
	}

	nonce, err := c.Provider.PendingNonceAt(c.ctx, txOp.From)
	if err != nil {
		logx.Errorf("pending nonce error %s", err.Error())
		return nil, err
	}

	gasPrice, err := c.Provider.SuggestGasPrice(c.ctx)
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
		signer := types.LatestSignerForChainID(c.chainId)
		txHashBytes := signer.Hash(tx).Bytes()
		signature, err := c.Kms.Sign(from, txHashBytes, "DIGEST")
		if err != nil {
			logx.Errorf("sign dynamic tx error %s", err.Error())
			return nil, err
		}
		return tx.WithSignature(signer, signature)
	}
	return txOp, nil
}

func (c *ChainClient) GetEstimateGas(tx *Tx, address string) (eGas uint64, err error) {
	toAddress := common.HexToAddress(tx.To)
	cMsg := ethereum.CallMsg{
		From: common.HexToAddress(address),
		To:   &toAddress,
		Data: tx.Data,
	}

	eGas, err = c.Provider.EstimateGas(c.ctx, cMsg)
	if err != nil {
		logx.Errorf("EstimateGas err:%s", err.Error())
		return 0, err
	}
	return eGas, nil
}

func (c *ChainClient) QueryRespection(hash common.Hash) bool {
	const retryTimes = 15
	var sucess = false
	for i := 0; i < retryTimes; i++ {
		receip, err := c.Provider.TransactionReceipt(c.ctx, hash)
		if err == nil {
			sucess = receip.BlockNumber.Int64() != 0
			break
		}
	}
	return sucess
}

func (c *ChainClient) watchToken(tokens []string, fromWallets []string) {
	currentBlock, err := c.Provider.BlockNumber(c.ctx)
	if err != nil {
		logx.Errorf("query block num error: %s", err.Error())
	}
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(currentBlock)),
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
		wowToken, err := wowfish.NewWowFishToken(tokenAddress, c.Provider)
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
	sub, err := c.Provider.SubscribeFilterLogs(c.ctx, query, logs)
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
				c.connectWsServer(c.ctx)
				//re watch
				sub, err = c.Provider.SubscribeFilterLogs(c.ctx, query, logs)
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

				from := common.HexToAddress(vLog.Topics[1].Hex()).String()
				to := common.HexToAddress(vLog.Topics[2].Hex()).String()

				//check callback address
				var valid = false
				for _, validAddress := range fromWallets {
					if strings.EqualFold(validAddress, from) || strings.EqualFold(validAddress, to) {
						valid = true
						break
					}
				}
				if !valid {
					continue
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

				//convert address to type
				if strings.EqualFold(from, c.bankAddress) {
					from = et.BankWallet
				} else {
					fromInfo, err := kms.NewKmsInstance().GetWalletInfo(from)
					if err == nil {
						from = fromInfo.Id.Type
					}
				}

				toInfo, err := kms.NewKmsInstance().GetWalletInfo(to)
				if err == nil { //"to address" maybe user
					to = toInfo.Id.Type
				}

				amount := contentValue.String()
				logx.Infof("From:%s,To:%s,Amount:%s",
					from,
					to,
					amount)
				//callback
				err = callback.Instance().Callback("/chain/consume", &callback.CallBackToWowfishData{
					CallbackBaseData: callback.CallbackBaseData{
						From: from,
						To:   to,
						Ret:  0,
					},
					Amount: amount,
				})
				if err != nil {
					logx.Errorf("Callback to game error %s", err.Error())
				}
			}
		}
	}(decimals)

}
