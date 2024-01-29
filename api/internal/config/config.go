package config

import "github.com/zeromicro/go-zero/rest"

type KmsConfig struct {
	Key           string
	Secret        string
	GameWallet    string
	DisposeWallet string
	GantWallet    string
}

type ChainConfig struct {
	Rpc      string
	WowToken string
	WoWNFT   string
	Callback string
}

type Config struct {
	rest.RestConf
	Chain     ChainConfig
	Kms       KmsConfig
	SecretKey string
}
