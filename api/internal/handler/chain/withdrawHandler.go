package chain

import (
	"net/http"

	"wowfish/api/internal/logic/chain"
	"wowfish/api/internal/svc"
	"wowfish/api/internal/types"
	"wowfish/api/pkg/response"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func WithdrawHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.WithdrawReqs
		if err := httpx.Parse(r, &req); err != nil {
			response.MakeError(r.Context(), w, err, response.ParamParseErrorCode)
			return
		}

		l := chain.NewWithdrawLogic(r.Context(), svcCtx)
		resp, err := l.Withdraw(&req)
		if err != nil {
			response.MakeError(r.Context(), w, err, response.LogicErroeCode)
		} else {
			response.MakeResponse(r.Context(), w, resp)
		}
	}
}
