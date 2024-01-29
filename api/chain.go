package main

import (
	"context"
	"flag"
	"os"

	"wowfish/api/internal/callback"
	"wowfish/api/internal/chain"
	"wowfish/api/internal/config"
	"wowfish/api/internal/handler"
	"wowfish/api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/wowfishconfig.yaml", "the config file")

func main() {
	flag.Parse()
	logModule := "file"
	if os.Getenv("DEBUG") == "true" {
		logModule = "console"
	}
	logConfig := logx.LogConf{
		Mode:     logModule,
		Rotation: "daily",
		KeepDays: 10,
		Path:     "./logs",
	}

	logx.MustSetup(logConfig)

	var c config.Config
	conf.MustLoad(*configFile, &c)

	//初始化chain client
	chainClient := chain.ChainClientInstance()
	err := chainClient.InitChainClient(context.Background(), &c)
	if err != nil {
		logx.Errorf("init chain client error %s", err.Error())
	}
	defer chainClient.Exit()

	callback.Instance().Init(c.Chain.Callback)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	logx.Infof("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
