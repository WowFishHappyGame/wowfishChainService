package chain

import (
	"net/http"

	"wowfish/api/internal/logic/chain"
	"wowfish/api/internal/svc"
	"wowfish/api/internal/types"
	"wowfish/api/pkg/response"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func MintNftHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.MintNftReqs
		if err := httpx.Parse(r, &req); err != nil {
			response.MakeError(r.Context(), w, err, response.ParamParseErrorCode)
			return
		}

		l := chain.NewMintNftLogic(r.Context(), svcCtx)
		resp, err := l.MintNft(&req)
		if err != nil {
			response.MakeError(r.Context(), w, err, response.ParamParseErrorCode)
		} else {
			response.MakeResponse(r.Context(), w, resp)
		}
	}
}
