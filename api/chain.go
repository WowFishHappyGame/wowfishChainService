package main

import (
	"context"
	"errors"
	"flag"
	"net/http"
	"os"

	"wowfish/api/internal/callback"
	"wowfish/api/internal/chain"
	"wowfish/api/internal/config"
	"wowfish/api/internal/handler"
	"wowfish/api/internal/svc"
	"wowfish/api/internal/util"
	"wowfish/api/pkg/response"
	"wowfish/api/pkg/ton"

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

	chainClient := chain.ChainClientInstance()
	err := chainClient.InitChainClient(context.Background(), &c)
	if err != nil {
		logx.Errorf("init chain client error %s", err.Error())
	}
	defer chainClient.Exit()

	err = ton.Instance().Init(c.Ton.Rpc)
	if err != nil {
		logx.Errorf("init ton rpc error %s", err)
	}

	callback.Instance().Init(c.Chain.Callback)
	util.Instance().Init(c.SecretKey)

	server := rest.MustNewServer(c.RestConf,
		rest.WithCustomCors( 
			func(header http.Header) {
				header.Set("Access-Control-Allow-Origin", "*")
				header.Set("Access-Control-Allow-Headers", "*")
				header.Set("Access-Control-Allow-Methods", " POST,OPTIONS")
			}, nil, "*"),
		rest.WithNotAllowedHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			response.MakeError(r.Context(), nil, errors.New("cors error"), response.NotAllowedError)
		})),
		rest.WithNotFoundHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			response.MakeError(r.Context(), nil, errors.New("cors error"), response.NotAllowedError)
		})))
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	logx.Infof("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
