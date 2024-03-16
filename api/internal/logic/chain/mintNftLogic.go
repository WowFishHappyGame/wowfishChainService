package chain

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	wowfish "wowfish/api/contract"
	"wowfish/api/internal/callback"
	"wowfish/api/internal/chain"
	"wowfish/api/internal/svc"
	"wowfish/api/internal/types"
	"wowfish/api/pkg/kms"
	"wowfish/api/pkg/response"
	et "wowfish/api/pkg/types"
	"wowfish/api/pkg/walletwatcher"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"
)

type MintNftLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMintNftLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MintNftLogic {
	return &MintNftLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MintNftLogic) MintNft(req *types.MintNftReqs) (resp *types.TransferResp, err error) {

	nft := l.svcCtx.Config.Chain.WoWNFT
	if nft == "" {
		return nil, errors.New("server error, nft config not find")
	}

	adminWallet := kms.NewKmsInstance().GetWalletAddress(et.AdminWallet)

	//async to mint
	go func(nftAddress string, adminWallet string, to string, id string) {

		walletwatcher.Instance().CheckAddressAmount(adminWallet)

		callback.Instance().DoWithCallback("/chain/mintnft",
			func() any {
				ret := callback.CallBackToWowfishNftData{
					CallbackBaseData: callback.CallbackBaseData{
						From: et.AdminWallet,
						To:   to,
						Ret:  0,
						Info: "",
					},
					Id: id,
				}
				chainMgr := chain.ChainClientInstance()
				token, err := wowfish.NewWowFishNft(common.HexToAddress(nftAddress), chainMgr.Provider)
				if err != nil {
					ret.Ret = response.TransferWowNftError
					ret.Info = err.Error()
					logx.Errorf("init nft error %s", err.Error())
					return ret
				}

				txOps, err := chainMgr.GenTranction(adminWallet)
				if err != nil {
					ret.Ret = response.TransferWowNftError
					ret.Info = err.Error()
					logx.Errorf("gen nft transcation op error %s ", err.Error())
					return ret
				}

				nftid := big.NewInt(0)
				nftid.SetString(id, 10)
				txs, err := token.SafeMint(txOps, common.HexToAddress(to), nftid, l.svcCtx.Config.Chain.NftUrl)
				if err != nil {
					ret.Ret = response.TransferWowNftError
					ret.Info = err.Error()
					logx.Errorf("transfer nft transcation op error %s ", err.Error())
					return ret
				}
				//nft not watch, so requse
				sucess := chainMgr.QueryRespection(txs.Hash())
				if !sucess {
					ret.Ret = response.TransferWowNftError
					errString := fmt.Sprintf("nft transcation on chain error %s ", txs.Hash())
					ret.Info = errString
					logx.Error(errString)
					return ret
				}
				ret.Info = txs.Hash().Hex()
				return ret
			})

	}(l.svcCtx.Config.Chain.WoWNFT, adminWallet, req.To, req.Id)
	resp = &types.TransferResp{
		Result: true,
	}
	return resp, nil
}
