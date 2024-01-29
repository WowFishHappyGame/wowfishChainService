package middleware

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"net/http"
	"net/url"
	"sort"
	"wowfish/api/pkg/response"

	"github.com/zeromicro/go-zero/core/logx"
)

type AuthMiddleware struct {
	secretKey string
}

func NewAuthMiddleware(secretKey string) *AuthMiddleware {
	return &AuthMiddleware{
		secretKey: secretKey,
	}
}

// GetSign 只支持value为字符串,只支持&连接
func getSign(params map[string]string, securityKey string) string {
	params["security_key"] = securityKey
	var keys []string
	for k := range params {
		if k != "sign" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	uParams := url.Values{}
	for _, k := range keys {
		uParams.Set(k, params[k])
	}
	data, _ := url.QueryUnescape(uParams.Encode())
	md5Str := toMD5(data)
	return md5Str
}
func toMD5(str string) string {
	hash := md5.New()
	hash.Write([]byte(str))
	return hex.EncodeToString(hash.Sum(nil))
}

func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO generate middleware implement function, delete after code implementation
		//签名校验
		const MaxMultipartMemory = 32 << 20 //32M
		if err := r.ParseMultipartForm(MaxMultipartMemory); err != nil {
			response.MakeError(r.Context(), w, err, response.NotAllowedError)
			return
		}
		values := r.PostForm
		params := make(map[string]string, 0)
		signIn := ""
		for key, val := range values {
			if len(val) > 0 && key != "sign" {
				params[key] = val[0]
			}
			if key == "sign" {
				signIn = val[0]
			}
		}

		sign := getSign(params, m.secretKey)
		if sign != signIn {
			logx.Errorf("sign error")
			response.MakeError(r.Context(), w, errors.New("sign error"), response.NotAllowedError)
			return
		}
		// Passthrough to next handler if need
		next(w, r)
	}
}
