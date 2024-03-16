package chainQuery

import (
	"context"

	"wowfish/api/internal/svc"
	"wowfish/api/internal/types"
	"wowfish/api/pkg/ton"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryTonWalletLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQueryTonWalletLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryTonWalletLogic {
	return &QueryTonWalletLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryTonWalletLogic) QueryTonWallet(req *types.TonWalletReq) (resp *types.TonWalletResp, err error) {
	// todo: add your logic here and delete this line
	wallet, err := l.svcCtx.Repo.GetWalletByDauthId(l.ctx, req.DAuthid)
	if err != nil {
		return nil, err
	}

	if wallet == "" {
		//create a new address
		tonWallet, err := ton.Instance().CreateWallet()
		if err != nil {
			return nil, err
		}

		err = l.svcCtx.Repo.InsertWallet(l.ctx, req.DAuthid, tonWallet.Address, tonWallet.Secret, tonWallet.Words)
		if err != nil {
			return nil, err
		}

		wallet = tonWallet.Address
	}
	return &types.TonWalletResp{
		Wallet: wallet,
	}, nil
}
