package callback

import (
	"errors"
	"strconv"
	"wowfish/api/pkg/dohttp"

	"github.com/zeromicro/go-zero/core/logx"
)

type CallBackToWowfish struct {
	callbackServer string
}

type CallBackToWowfishData struct {
	From   string `from:"form"`
	To     string `from:"to"`
	Amount string `from:"amount"`
	Ret    int64  `from:"ret"`
}

var ins = CallBackToWowfish{
	callbackServer: "",
}

func (this *CallBackToWowfish) Init(callbackUrl string) {
	this.callbackServer = callbackUrl
}

func Instance() *CallBackToWowfish {
	return &ins
}

func (this *CallBackToWowfish) Callback(data *CallBackToWowfishData) error {
	if this.callbackServer != "" {
		rsp, err := dohttp.DoMultiFormHttp(map[string]string{}, "POST", this.callbackServer,
			map[string]string{
				"ret":     strconv.FormatInt(data.Ret, 10),
				"address": data.From,
				"to":      data.To,
				"amount":  data.Amount,
			})
		defer rsp.Body.Close()
		if nil != err {
			logx.Errorf("post to callback error %s", err.Error())
			return err
		}
	}
	logx.Errorf("callback data is:%+v", data)
	return errors.New("callback url is nui")
}
