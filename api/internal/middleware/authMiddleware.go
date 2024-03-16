package middleware

import (
	"errors"
	"net/http"
	"wowfish/api/internal/util"
	"wowfish/api/pkg/response"

	"github.com/zeromicro/go-zero/core/logx"
)

type AuthMiddleware struct {
}

func NewAuthMiddleware(secretKey string) *AuthMiddleware {
	return &AuthMiddleware{}
}

func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO generate middleware implement function, delete after code implementation
		const MaxMultipartMemory = 32 << 20 //32M
		if err := r.ParseMultipartForm(MaxMultipartMemory); err != nil {
			response.MakeError(r.Context(), w, err, response.NotAllowedError)
			return
		}
		values := r.PostForm
		params := make(map[string]any, 0)
		signIn := ""
		for key, val := range values {
			if len(val) > 0 && key != "sign" {
				params[key] = val[0]
			}
			if key == "sign" {
				signIn = val[0]
			}
		}

		sign := util.Instance().GetInfrasSign(params)
		if sign != signIn {
			logx.Errorf("sign error")
			response.MakeError(r.Context(), w, errors.New("sign error"), response.NotAllowedError)
			return
		}
		// Passthrough to next handler if need
		next(w, r)
	}
}
