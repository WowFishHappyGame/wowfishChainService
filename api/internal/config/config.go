package config

import "github.com/zeromicro/go-zero/rest"

type KmsConfig struct {
	Key           string
	Secret        string
	GameWallet    string
	DisposeWallet string
	GainWallet    string
	AdminWallet   string
}

type ChainConfig struct {
	Rpc      string
	WowToken string
	WoWNFT   string
	GameBank string
	Callback string
	NftUrl   string
}

type Metrics struct {
	EMailServer string
	Threshold   float64
	ToEmail     string
	SignKey     string
}

type Config struct {
	rest.RestConf
	Chain     ChainConfig
	Kms       KmsConfig
	SecretKey string
	Metrics   Metrics
	Mysql     Mysql
	Ton       Ton
}

type Mysql struct {
	DataSource string
}

type Ton struct {
	Rpc string
}
