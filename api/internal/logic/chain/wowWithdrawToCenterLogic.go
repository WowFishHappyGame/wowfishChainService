package chain

import (
	"context"
	"fmt"

	wowfish "wowfish/api/contract"
	"wowfish/api/internal/callback"
	"wowfish/api/internal/chain"
	"wowfish/api/internal/svc"
	"wowfish/api/internal/types"
	"wowfish/api/pkg/kms"
	"wowfish/api/pkg/response"
	"wowfish/api/pkg/walletwatcher"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"

	et "wowfish/api/pkg/types"
)

type WowWithdrawToCenterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWowWithdrawToCenterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WowWithdrawToCenterLogic {
	return &WowWithdrawToCenterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WowWithdrawToCenterLogic) WowWithdrawToCenter(req *types.WithdrawReqs) (resp *types.TransferResp, err error) {
	//async to do transfer

	bankAddress := l.svcCtx.Config.Chain.GameBank

	centerAddress := kms.NewKmsInstance().GetWalletAddress(et.GameWallet)
	if centerAddress == "" {
		logx.Errorf("query gamewallet is null")
		return nil, fmt.Errorf("query gamewallet is null")
	}

	adminAddress := kms.NewKmsInstance().GetWalletAddress(et.AdminWallet)
	if centerAddress == "" {
		logx.Errorf("query adminwallet is null")
		return nil, fmt.Errorf("query adminwallet is null")
	}

	go func(bankAddress string, centerAddress string, adminAddress string) {
		walletwatcher.Instance().CheckAddressAmount(adminAddress)

		callback.Instance().DoWithCallback("/chain/consume", func() any {
			ret := &callback.CallBackToWowfishData{
				CallbackBaseData: callback.CallbackBaseData{
					From: et.AdminWallet,
					To:   et.GameWallet,
					Ret:  0,
					Info: "",
				},
				Amount: "",
			}
			chainMgr := chain.ChainClientInstance()
			wowfishBank, err := wowfish.NewWowFishBank(common.HexToAddress(bankAddress), chainMgr.Provider)
			if err != nil {
				ret.Ret = response.TransferWowTokenError
				ret.Info = err.Error()
				return ret
			}
			txOp, err := chainMgr.GenTranction(adminAddress)
			if err != nil {
				ret.Ret = response.TransferWowTokenError
				ret.Info = err.Error()
				return ret
			}
			tx, err := wowfishBank.Withdraw(txOp, common.HexToAddress(centerAddress))
			if err != nil {
				logx.Errorf("withdraw error %s", err.Error())
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
	}(bankAddress, centerAddress, adminAddress)

	resp = &types.TransferResp{
		Result: true,
	}
	return resp, nil
}
