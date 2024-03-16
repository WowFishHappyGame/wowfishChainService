package util

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/url"
	"reflect"
	"sort"

	"github.com/zeromicro/go-zero/core/logx"
)

type Signer struct {
	Key string
}

var ins = Signer{
	Key: "",
}

func (s *Signer) Init(key string) {
	s.Key = key
}

func Instance() *Signer {
	return &ins
}

func (s *Signer) EncodeSignData(params map[string]any) string {
	var keys []string
	for k := range params {
		if k != "sign" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	uParams := url.Values{}
	logx.Infof("begin to sign %+v", params)
	for _, k := range keys {
		v := reflect.ValueOf(params[k])
		var param string
		switch v.Kind() {
		case reflect.String:
			param = v.String()
		case reflect.Int64:
			param = fmt.Sprintf("%d", v.Int())
		case reflect.Float64:
			param = fmt.Sprintf("%d", int64(v.Float()))
		default:
			logx.Errorf("param type error, %s -- %v  is not invaliable", k, v.Kind())
			continue
		}
		logx.Info(param)
		uParams.Set(k, param)
	}
	data, _ := url.QueryUnescape(uParams.Encode())
	return data
}


func (s *Signer) GetInfrasSign(params map[string]any) string {
	params["security_key"] = s.Key
	data := s.EncodeSignData(params)
	md5Str := s.ToMD5(data)
	return md5Str
}
func (s *Signer) ToMD5(str string) string {
	hash := md5.New()
	_, err := hash.Write([]byte(str))
	if err != nil {
		logx.Errorf("%s to md5 error %s", str, err.Error())
		return ""
	}
	return hex.EncodeToString(hash.Sum(nil))
}
