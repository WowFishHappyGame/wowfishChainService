package chain

import (
	"context"
	"encoding/json"
	"math"
	"math/big"

	wowfish "wowfish/api/contract"
	"wowfish/api/internal/svc"
	"wowfish/api/internal/types"
	"wowfish/api/pkg/chain"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"
)

type WithdrawReq struct {
	Address string  `json:"address"`
	Amount  float64 `json:"amount"`
}

type WithdrawLogic struct {
	logx.Logger
	ctx          context.Context
	svcCtx       *svc.ServiceContext
	wowfishToken *wowfish.WowFishToken
	decimal      uint8
}

func NewWithdrawLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WithdrawLogic {
	tAddress := common.HexToAddress(svcCtx.Config.Chain.WowToken)
	chainMgr := chain.ChainClientInstance()
	wowfishToken, err := wowfish.NewWowFishToken(tAddress, chainMgr.Provider)
	if err != nil {
		logx.Errorf("new erc20 token error %s", err.Error())
		return nil
	}

	decimal, err := wowfishToken.Decimals(&bind.CallOpts{})
	if err != nil {
		logx.Errorf("new erc20 token error %s", err.Error())
		return nil
	}

	return &WithdrawLogic{
		Logger:       logx.WithContext(ctx),
		ctx:          ctx,
		svcCtx:       svcCtx,
		wowfishToken: wowfishToken,
		decimal:      decimal,
	}
}

func (l *WithdrawLogic) Withdraw(req *types.WithdrawReqs) (resp *types.WithdrawResp, err error) {
	// todo: add your logic here and delete this line
	//construct transfer data
	var data []wowfish.TransferData
	var reqs []WithdrawReq
	err = json.Unmarshal([]byte(req.Reqs), &reqs)
	if err != nil {
		logx.Errorf("Unmarshal req %s", err.Error())
		return
	}
	for _, value := range reqs {
		data = append(data, wowfish.TransferData{
			To:     common.HexToAddress(value.Address),
			Amount: big.NewInt(int64(value.Amount * math.Pow10(int(l.decimal)))),
		})
	}
	provider := chain.ChainClientInstance()
	abi, err := wowfish.WowFishTokenMetaData.GetAbi()
	if err != nil {
		logx.Errorf("get token abi error %s", err.Error())
		return
	}
	txData, err := abi.Pack("transferMulti", data)
	if err != nil {
		logx.Errorf("abi pack TransferMulti error %s", err.Error())
		return
	}
	txHash, err := provider.CommitTranscation(&chain.Tx{
		To:    l.svcCtx.Config.Chain.WowToken,
		Value: 0,
		Data:  txData,
	})
	if err != nil {
		logx.Errorf("TransferMulti error %s", err.Error())
		return
	}
	logx.Infof("transcation res %s", txHash)
	resp = &types.WithdrawResp{
		Result: true,
	}
	return
}
