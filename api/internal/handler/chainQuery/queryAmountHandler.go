package chainQuery

import (
	"net/http"

	"wowfish/api/internal/logic/chainQuery"
	"wowfish/api/internal/svc"
	"wowfish/api/internal/types"
	"wowfish/api/pkg/response"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func QueryAmountHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.WoWQueryAmountReq
		if err := httpx.Parse(r, &req); err != nil {
			response.MakeError(r.Context(), w, err, response.ParamParseErrorCode)
		}
		l := chainQuery.NewQueryAmountLogic(r.Context(), svcCtx)
		resp, err := l.QueryAmount(&req)
		if err != nil {
			response.MakeError(r.Context(), w, err, response.ParamParseErrorCode)
		} else {
			response.MakeResponse(r.Context(), w, resp)
		}
	}
}
