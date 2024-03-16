package chainQuery

import (
	"context"
	"fmt"
	"math"
	"math/big"

	wowfish "wowfish/api/contract"
	"wowfish/api/internal/chain"
	"wowfish/api/internal/svc"
	"wowfish/api/internal/types"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"

	"wowfish/api/pkg/kms"
	et "wowfish/api/pkg/types"
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

func (l *QueryAmountLogic) QueryAmount(req *types.WalletType) (resp *types.WoWQueryAmountResp, err error) {
	// todo: add your logic here and delete this line
	// todo: add your logic here and delete this line
	tAddress := common.HexToAddress(l.svcCtx.Config.Chain.WowToken)
	chainMgr := chain.ChainClientInstance()
	var queryAddress string
	if req.Type == et.BankWallet {
		queryAddress = l.svcCtx.Config.Chain.GameBank
	} else {
		queryAddress = kms.NewKmsInstance().GetWalletAddress(req.Type)
	}
	if queryAddress == "" {
		logx.Errorf("(%#v) QueryAmount erc20 token error not support type %s", req, req.Type)
		return nil, fmt.Errorf("QueryAmount erc20 token error not support type %s", req.Type)
	}
	wowfishToken, err := wowfish.NewWowFishToken(tAddress, chainMgr.Provider)
	if err != nil {
		logx.Errorf("(%#v) QueryAmount new erc20 token error %s", req, err.Error())
		return nil, err
	}

	decimal, err := wowfishToken.Decimals(&bind.CallOpts{})
	if err != nil {
		logx.Errorf("(%#v) QueryAmount rc20 Decimals error %s", req, err.Error())
		return nil, err
	}

	add := common.HexToAddress(queryAddress)
	amount, err := wowfishToken.BalanceOf(&bind.CallOpts{}, add)
	if err != nil {
		logx.Errorf("(%#v) QueryAmount rc20 BalanceOf error %s", req, err.Error())
		return nil, err
	}
	a := big.NewFloat(float64(amount.Int64()))
	b := big.NewFloat(math.Pow10(int(decimal)))

	quo := big.NewFloat(0)
	quo = quo.Quo(a, b)
	return &types.WoWQueryAmountResp{
		Amount: quo.String(),
	}, nil
}
