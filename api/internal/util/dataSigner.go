package util

import (
	"crypto/md5"
	"encoding/hex"
	"net/url"
	"sort"
)

type Signer struct {
	Key string
}

var ins = Signer{
	Key: "",
}

func (this *Signer) Init(key string) {
	this.Key = key
}

func Instance() *Signer {
	return &ins
}

// GetSign 只支持value为字符串,只支持&连接
func (this *Signer) GetSign(params map[string]string) string {
	params["security_key"] = this.Key
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
