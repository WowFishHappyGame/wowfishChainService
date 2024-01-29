package chainQuery

import (
	"context"
	"math"
	"math/big"

	wowfish "wowfish/api/contract"
	"wowfish/api/internal/chain"
	"wowfish/api/internal/svc"
	"wowfish/api/internal/types"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"
)

type QueryAmountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQueryAmountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryAmountLogic {
	return &QueryAmountLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryAmountLogic) QueryAmount(req *types.WoWQueryAmountReq) (resp float64, err error) {
	// todo: add your logic here and delete this line
	tAddress := common.HexToAddress(l.svcCtx.Config.Chain.WowToken)
	chainMgr := chain.ChainClientInstance()
	wowfishToken, err := wowfish.NewWowFishToken(tAddress, chainMgr.Provider)
	if err != nil {
		logx.Errorf("QueryAmount new erc20 token error %s", err.Error())
		return 0, err
	}

	decimal, err := wowfishToken.Decimals(&bind.CallOpts{})
	if err != nil {
		logx.Errorf("QueryAmount rc20 Decimals error %s", err.Error())
		return 0, err
	}

	add := common.HexToAddress(req.Address)
	amount, err := wowfishToken.BalanceOf(&bind.CallOpts{}, add)
	if err != nil {
		logx.Errorf("QueryAmount rc20 BalanceOf error %s", err.Error())
		return 0, err
	}
	a := big.NewFloat(float64(amount.Int64()))
	b := big.NewFloat(math.Pow10(int(decimal)))

	quo := big.NewFloat(0).Quo(a, b)
	ret, _ := quo.Float64()
	return ret, nil
}
