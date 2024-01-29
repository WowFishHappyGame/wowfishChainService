package chain

import (
	"context"
	"errors"
	"math"
	"math/big"
	"strconv"

	wowfish "wowfish/api/contract"
	"wowfish/api/internal/callback"
	"wowfish/api/internal/chain"
	"wowfish/api/internal/svc"
	"wowfish/api/internal/types"
	"wowfish/api/pkg/response"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"
)

type WowTransferLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWowTransferLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WowTransferLogic {
	return &WowTransferLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WowTransferLogic) commitTransToChain(from string, to string, amount string) {

	ret := callback.CallBackToWowfishData{
		From:   from,
		To:     to,
		Amount: amount,
		Ret:    0,
	}

	defer func() {
		if ret.Ret != 0 {
			callback.Instance().Callback(&ret)
		}
	}()

	tAddress := common.HexToAddress(l.svcCtx.Config.Chain.WowToken)
	chainMgr := chain.ChainClientInstance()
	wowfishToken, err := wowfish.NewWowFishToken(tAddress, chainMgr.Provider)
	if err != nil {
		logx.Errorf("new erc20 token error %s", err.Error())
		ret.Ret = response.TransferWowTokenError
		return
	}

	decimal, err := wowfishToken.Decimals(&bind.CallOpts{})
	if err != nil {
		logx.Errorf("new erc20 token error %s", err.Error())
		ret.Ret = response.TransferWowTokenError
		return
	}

	abi, err := wowfish.WowFishTokenMetaData.GetAbi()
	if err != nil {
		logx.Errorf("get token abi error %s", err.Error())
		ret.Ret = response.TransferWowTokenError
		return
	}
	fv, err := strconv.ParseFloat(amount, 64)

	txData, err := abi.Pack("transfer", common.HexToAddress(to), big.NewInt(int64(fv*math.Pow10(int(decimal)))))
	if err != nil {
		logx.Errorf("abi pack transfer error %s", err.Error())
		ret.Ret = response.TransferWowTokenError
		return
	}
	txHash, err := chainMgr.CommitTranscation(&chain.Tx{
		To:    l.svcCtx.Config.Chain.WowToken,
		Value: 0,
		Data:  txData,
	}, from)
	if err != nil {
		logx.Errorf("transfer error %s", err.Error())
		ret.Ret = response.TransferWowTokenError
		return
	}
	logx.Infof("transcation res %s", txHash)
}

func (l *WowTransferLogic) WowTransfer(req *types.TransferReqs) (resp *types.TransferResp, err error) {
	//coroutine the rpc may timeout
	if req.From == "" || req.To == "" || req.Amount == "" || req.Amount == "0" {
		return nil, errors.New("param error")
	}
	go l.commitTransToChain(req.From, req.To, req.Amount)
	resp = &types.TransferResp{
		Result: true,
	}
	return
}
