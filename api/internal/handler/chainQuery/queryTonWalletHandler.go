package chainQuery

import (
	"net/http"

	"wowfish/api/internal/logic/chainQuery"
	"wowfish/api/internal/svc"
	"wowfish/api/internal/types"
	"wowfish/api/pkg/response"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func QueryTonWalletHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.TonWalletReq
		if err := httpx.Parse(r, &req); err != nil {
			response.MakeError(r.Context(), w, err, response.ParamParseErrorCode)
			return
		}

		l := chainQuery.NewQueryTonWalletLogic(r.Context(), svcCtx)
		resp, err := l.QueryTonWallet(&req)
		if err != nil {
			response.MakeError(r.Context(), w, err, response.LogicErroeCode)
		} else {
			response.MakeResponse(r.Context(), w, resp)
		}
	}
}
