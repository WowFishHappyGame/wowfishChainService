package svc

import (
	"wowfish/api/internal/config"
	"wowfish/api/internal/middleware"

	"github.com/zeromicro/go-zero/rest"

	"wowfish/api/internal/db"
)

type ServiceContext struct {
	Config         config.Config
	AuthMiddleware rest.Middleware
	Repo           db.TonRepo
}

func NewServiceContext(c config.Config) *ServiceContext {

	return &ServiceContext{
		Config:         c,
		AuthMiddleware: middleware.NewAuthMiddleware(c.SecretKey).Handle,
		Repo:           db.NewTonRepo(c),
	}
}
