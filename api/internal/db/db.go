package db

import (
	"context"
	"wowfish/api/internal/config"
	"wowfish/api/internal/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type TonRepo interface {
	GetWalletByDauthId(ctx context.Context, dauthid string) (string, error)
	InsertWallet(ctx context.Context, dauthid string, address string, secret string, words string) error
}

type TonRepoImpl struct {
	model.TWalletAddressModel
}

func NewTonRepo(c config.Config) TonRepo {
	dbConn := sqlx.NewMysql(c.Mysql.DataSource)
	return TonRepoImpl{
		TWalletAddressModel: model.NewTWalletAddressModel(dbConn),
	}
}
