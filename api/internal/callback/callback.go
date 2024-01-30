package callback

import (
	"encoding/json"
	"errors"
	"strconv"
	"wowfish/api/internal/util"
	"wowfish/api/pkg/dohttp"

	"github.com/zeromicro/go-zero/core/logx"
)

type CallBackToWowfish struct {
	callbackServer string
}

type CallBackToWowfishData struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount string `json:"amount"`
	Ret    int64  `json:"ret"`
	Info   string `json:"info"`
	Sign   string `json:"sign"`
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

		originData := map[string]string{
			"ret":    strconv.FormatInt(data.Ret, 10),
			"from":   data.From,
			"to":     data.To,
			"amount": data.Amount,
			"info":   data.Info,
		}

		sign := util.Instance().GetInfrasSign(originData)
		data.Sign = sign

		jsonData, err := json.Marshal(data)

		rsp, err := dohttp.DoJsonHttp(map[string]string{}, "POST", this.callbackServer, jsonData)
		defer rsp.Body.Close()
		if nil != err {
			logx.Errorf("post to callback error %s", err.Error())
			return err
		}
	}
	logx.Errorf("callback data is:%+v", data)
	return errors.New("callback url is nui")
}
