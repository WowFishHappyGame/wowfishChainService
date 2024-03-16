package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var (
	_ TWalletAddressModel = (*customTWalletAddressModel)(nil)
)

type (
	// TWalletAddressModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTWalletAddressModel.
	TWalletAddressModel interface {
		tWalletAddressModel
		GetWalletByDauthId(ctx context.Context, dauthid string) (string, error)
		InsertWallet(ctx context.Context, dauthid string, address string, secret string, words string) error
	}

	customTWalletAddressModel struct {
		*defaultTWalletAddressModel
	}
)

// NewTWalletAddressModel returns a model for the database table.
func NewTWalletAddressModel(conn sqlx.SqlConn) TWalletAddressModel {
	return &customTWalletAddressModel{
		defaultTWalletAddressModel: newTWalletAddressModel(conn),
	}
}

func (m *customTWalletAddressModel) GetWalletByDauthId(ctx context.Context, dauthid string) (string, error) {
	query := fmt.Sprintf("select `address` from %s where `dauth_id`=?", m.table)
	var resp string
	err := m.conn.QueryRowPartialCtx(ctx, &resp, query, dauthid)
	switch err {
	case nil:
		return resp, nil
	case ErrNotFound:
		return "", nil
	default:
		return "", err
	}
}

func (m *customTWalletAddressModel) InsertWallet(ctx context.Context, dauthid string, address string, secret string, words string) error {
	modle := &TWalletAddress{
		Address: address,
		Secret:  secret,
		Words:   words,
		DauthId: dauthid,
	}
	_, err := m.Insert(ctx, modle)
	return err
}
