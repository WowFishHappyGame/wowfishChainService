package ton

import (
	"context"
	"encoding/hex"
	"strings"
	"sync"

	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/wallet"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetWalletInfo struct {
	Address string
	Secret  string
	Words   string
}

type Ton struct {
	rpcUrl string
	api    *ton.APIClient
}

var once sync.Once
var instance Ton

func Instance() *Ton {
	once.Do(func() {
		instance = Ton{}
	})
	return &instance
}

func (t *Ton) Init(rpc string) error {
	t.rpcUrl = rpc
	client := liteclient.NewConnectionPool()

	// get config
	cfg, err := liteclient.GetConfigFromUrl(context.Background(), t.rpcUrl)
	if err != nil {
		logx.Errorf("get config err: %s", err.Error())
		return err
	}

	// connect to mainnet lite servers
	err = client.AddConnectionsFromConfig(context.Background(), cfg)
	if err != nil {
		logx.Errorf("connection err: %s", err.Error())
		return err
	}

	// api client with full proof checks
	t.api = ton.NewAPIClient(client, ton.ProofCheckPolicySecure)
	t.api.WithRetry()
	t.api.SetTrustedBlockFromConfig(cfg)
	return nil
}

func (t *Ton) CreateWallet() (*GetWalletInfo, error) {

	// bound all requests to single ton node
	words := wallet.NewSeed()
	w, err := wallet.FromSeed(t.api, words, wallet.V4R2)
	if err != nil {
		logx.Errorf("FromSeed err: %s", err.Error())
		return nil, err
	}

	return &GetWalletInfo{
		Address: w.WalletAddress().String(),
		Words:   strings.Join(words, ","),
		Secret:  hex.EncodeToString(w.PrivateKey()),
	}, nil
}
