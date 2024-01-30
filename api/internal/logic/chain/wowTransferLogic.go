package chain

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"math"
	"math/big"
	"strconv"

	wowfish "wowfish/api/contract"
	"wowfish/api/internal/callback"
	"wowfish/api/internal/chain"
	"wowfish/api/internal/svc"
	"wowfish/api/internal/types"
	"wowfish/api/internal/util"
	"wowfish/api/pkg/dohttp"
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

func (l *WowTransferLogic) checkAddressAmount(address string) {
	if l.svcCtx.Config.Metrics.EMailServer != "" {
		chainMgr := chain.ChainClientInstance()
		balance, err := chainMgr.Provider.BalanceAt(context.Background(), common.HexToAddress(address), nil)
		if err != nil {
			logx.Error("query blance error ", err.Error())
			return
		}
		balanceInt := balance.Int64()
		if balanceInt < int64(l.svcCtx.Config.Metrics.Threshold*math.Pow10(18)) {
			//send the warning email
			emailData := map[string]string{
				"email":   l.svcCtx.Config.Metrics.ToEmail,
				"name":    "Wowfish",
				"subject": "Wallet Balance not enough",
				"body":    address + ":" + strconv.FormatFloat(float64(balanceInt)/math.Pow10(18), 'f', 10, 64),
			}
			encodeData := util.Instance().EncodeSignData(emailData)
			encodeData += "&" + l.svcCtx.Config.Metrics.SignKey
			sign := util.Instance().ToMD5(encodeData)

			emailData["sign"] = sign

			rsp, err := dohttp.DoMultiFormHttp(map[string]string{}, "POST",
				l.svcCtx.Config.Metrics.EMailServer,
				emailData)
			if err != nil {
				logx.Error("Send Metricsati email error", err.Error())
				return
			}
			defer rsp.Body.Close()

			body, err := io.ReadAll(rsp.Body)
			if err != nil {
				logx.Error("Send Metricsati email read response", err.Error())
				return
			}

			type EmailResponse struct {
				Code int64 `json:"code"`
			}

			var b = EmailResponse{}

			err = json.Unmarshal(body, &b)
			if err != nil {
				logx.Error("Send Metricsati Unmarshal response", err.Error())
				return
			}
			if b.Code != 1 {
				logx.Error("Send Metricsati email errcode", b.Code)
				return
			}
			logx.Info("send metrics email sucess")
		}
	}
}

func (l *WowTransferLogic) commitTransToChain(from string, to string, amount string) {

	l.checkAddressAmount(from)

	ret := callback.CallBackToWowfishData{
		From:   from,
		To:     to,
		Amount: amount,
		Ret:    0,
		Info:   "",
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
		ret.Info = err.Error()
		return
	}

	decimal, err := wowfishToken.Decimals(&bind.CallOpts{})
	if err != nil {
		logx.Errorf("new erc20 token error %s", err.Error())
		ret.Ret = response.TransferWowTokenError
		ret.Info = err.Error()
		return
	}

	abi, err := wowfish.WowFishTokenMetaData.GetAbi()
	if err != nil {
		logx.Errorf("get token abi error %s", err.Error())
		ret.Ret = response.TransferWowTokenError
		ret.Info = err.Error()
		return
	}
	fv, err := strconv.ParseFloat(amount, 64)

	txData, err := abi.Pack("transfer", common.HexToAddress(to), big.NewInt(int64(fv*math.Pow10(int(decimal)))))
	if err != nil {
		logx.Errorf("abi pack transfer error %s", err.Error())
		ret.Ret = response.TransferWowTokenError
		ret.Info = err.Error()
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
		ret.Info = err.Error()
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
