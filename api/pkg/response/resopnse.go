package response

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

type NormalJsonRes struct {
	Ret  int64  `json:"ret"`
	Info string `json:"info"`
	Data any    `json:"data"`
}

func MakeResponse(ctx context.Context, w http.ResponseWriter, v any) {

	normalRet := &NormalJsonRes{}
	normalRet.Ret = 0
	normalRet.Info = "ok"
	normalRet.Data = v
	httpx.OkJsonCtx(ctx, w, normalRet)
}

func MakeError(ctx context.Context, w http.ResponseWriter, err error, errorCode int64,
	fns ...func(w http.ResponseWriter, err error)) {

	normalRet := &NormalJsonRes{}
	normalRet.Ret = errorCode
	normalRet.Info = err.Error()
	normalRet.Data = nil
	httpx.OkJsonCtx(ctx, w, normalRet)
}
