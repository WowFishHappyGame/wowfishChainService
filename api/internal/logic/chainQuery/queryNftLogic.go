package chainQuery

import (
	"context"
	"math/big"
	wowfish "wowfish/api/contract"
	"wowfish/api/internal/chain"
	"wowfish/api/internal/svc"
	"wowfish/api/internal/types"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"
)

type QueryNftLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQueryNftLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryNftLogic {
	return &QueryNftLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryNftLogic) QueryNft(req *types.NftQueryIdsReq) (resp *types.NftQueryIdsResp, err error) {
	// todo: add your logic here and delete this line
	nftAddress := common.HexToAddress(l.svcCtx.Config.Chain.WoWNFT)
	chainMgr := chain.ChainClientInstance()
	nftToken, err := wowfish.NewWowFishNft(nftAddress, chainMgr.Provider)
	if err != nil {
		logx.Errorf("Nft token error %s", err.Error())
		return nil, err
	}

	ownerAddress := common.HexToAddress(req.Address)
	count, err := nftToken.BalanceOf(&bind.CallOpts{}, ownerAddress)
	if err != nil {
		logx.Errorf("get nft count error %s", err.Error())
		return nil, err
	}
	resp = &types.NftQueryIdsResp{
		Ids: []int64{},
	}

	var index int64
	for index = 0; index < count.Int64(); index++ {
		id, err := nftToken.TokenOfOwnerByIndex(&bind.CallOpts{}, ownerAddress, big.NewInt(index))
		if err != nil {
			logx.Errorf("get nft by index error %s", err.Error())
			return nil, err
		}
		resp.Ids = append(resp.Ids, id.Int64())
	}
	return resp, nil
}
