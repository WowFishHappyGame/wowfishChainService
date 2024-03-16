package chain

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/big"

	"wowfish/api/internal/callback"
	"wowfish/api/internal/chain"
	"wowfish/api/internal/svc"
	"wowfish/api/internal/types"
	"wowfish/api/pkg/kms"
	"wowfish/api/pkg/response"
	"wowfish/api/pkg/walletwatcher"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"

	et "wowfish/api/pkg/types"

	wowfish "wowfish/api/contract"
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

func (l *WowTransferLogic) commitTransToChain(walletType string, from string, to string, amount string) {

	walletwatcher.Instance().CheckAddressAmount(from)

	callback.Instance().DoWithCallback("/chain/consume", func() any {

		ret := &callback.CallBackToWowfishData{
			CallbackBaseData: callback.CallbackBaseData{
				From: walletType,
				To:   to,
				Ret:  0,
				Info: "",
			},
			Amount: amount,
		}
		tAddress := common.HexToAddress(l.svcCtx.Config.Chain.WowToken)
		chainMgr := chain.ChainClientInstance()

		wowfishToken, err := wowfish.NewWowFishToken(tAddress, chainMgr.Provider)
		if err != nil {
			logx.Errorf("new erc20 token error %s", err.Error())
			ret.Ret = response.TransferWowTokenError
			ret.Info = err.Error()
			return ret
		}

		decimal, err := wowfishToken.Decimals(&bind.CallOpts{})
		if err != nil {
			logx.Errorf("new erc20 token error %s", err.Error())
			ret.Ret = response.TransferWowTokenError
			ret.Info = err.Error()
			return ret
		}

		txOp, err := chainMgr.GenTranction(from)
		if err != nil {
			logx.Errorf("gen transcation op error %s", err.Error())
			ret.Ret = response.TransferWowTokenError
			ret.Info = err.Error()
			return ret
		}

		fv, _, err := big.ParseFloat(amount, 10, 256, big.ToNearestEven)
		if err != nil {
			logx.Errorf("parse float error %s", err.Error())
			ret.Ret = response.TransferWowTokenError
			ret.Info = err.Error()
			return ret
		}

		var amountInt64 = big.NewInt(0)
		fv.Mul(fv, big.NewFloat(math.Pow10(int(decimal)))).Int(amountInt64)
		tx, err := wowfishToken.Transfer(txOp, common.HexToAddress(to), amountInt64)
		if err != nil {
			logx.Errorf("transfer error %s", err.Error())
			ret.Ret = response.TransferWowTokenError
			ret.Info = err.Error()
			return ret
		}
		logx.Infof("transcation res %s", tx.Hash().String())
		//check trans recipect
		sucess := chainMgr.QueryRespection(tx.Hash())
		if !sucess {
			ret.Ret = response.TransferWowTokenError
			ret.Info = fmt.Sprintf("Transfer wow token on chain error %s", tx.Hash().Hex())
			return ret
		}
		return nil
	})
}

func (l *WowTransferLogic) WowTransfer(req *types.TransferReqs) (resp *types.TransferResp, err error) {
	//coroutine the rpc may timeout
	if req.Type == "" || req.To == "" || req.Amount == "" || req.Amount == "0" {
		return nil, errors.New("param error")
	}
	if req.Type != et.DisposeWallet &&
		req.Type != et.GameWallet {
		return nil, errors.New("only support from DisposeWallet or GameWallet")
	}
	from := kms.NewKmsInstance().GetWalletAddress(req.Type)
	if from == "" {
		return nil, fmt.Errorf("can't find wallet witch type named %s", req.Type)
	}
	to := req.To
	if to == et.BankWallet || to == et.AdminWallet {
		return nil, fmt.Errorf("can't transfer to wallet %s", to)
	}

	if to == et.DisposeWallet || to == et.GainWallet ||
		to == et.GameWallet {
		to = kms.NewKmsInstance().GetWalletAddress(to)
	}
	if to == "" {
		return nil, fmt.Errorf("can't find wallet witch type named %s", to)
	}

	go l.commitTransToChain(req.Type, from, to, req.Amount)
	resp = &types.TransferResp{
		Result: true,
	}
	return
}
