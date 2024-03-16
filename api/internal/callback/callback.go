package callback

import (
	"encoding/json"
	"errors"
	"wowfish/api/internal/util"
	"wowfish/api/pkg/dohttp"

	"github.com/zeromicro/go-zero/core/logx"
)

type CallBackToWowfish struct {
	callbackServer string
}

type CallbackLogicFunc func() any

type CallbackBaseData struct {
	From string `json:"from"`
	To   string `json:"to"`
	Ret  int64  `json:"ret"`
	Info string `json:"info"`
	Sign string `json:"sign"`
}

type CallBackToWowfishData struct {
	CallbackBaseData
	Amount string `json:"amount"`
}

type CallBackToWowfishNftData struct {
	CallbackBaseData
	Id string `json:"nft_id"`
}

var ins = CallBackToWowfish{
	callbackServer: "",
}

func (c *CallBackToWowfish) Init(callbackUrl string) {
	c.callbackServer = callbackUrl
}

func Instance() *CallBackToWowfish {
	return &ins
}

func (c *CallBackToWowfish) DoWithCallback(router string, fun CallbackLogicFunc) {
	data := fun()
	if data != nil { //if nil  callback process to watcher
		err := c.Callback(router, data)
		if err != nil {
			logx.Errorf("callback to game error %s", err.Error())
		}
	}
}

func (c *CallBackToWowfish) Callback(router string, data any) error {
	if c.callbackServer != "" {

		jsonData, err := json.Marshal(data)
		if err != nil {
			logx.Errorf("marshal data to json error %s", err.Error())
		}

		originData := make(map[string]any)
		err = json.Unmarshal(jsonData, &originData)
		if err != nil {
			logx.Errorf("unmarshal json to map error %s", err.Error())
		}

		sign := util.Instance().GetInfrasSign(originData)
		originData["sign"] = sign
		delete(originData, "security_key")

		//converto json again

		sendData, err := json.Marshal(originData)
		if err != nil {
			logx.Errorf("marshal data to json error %s", err.Error())
		}

		logx.Infof("send to callback %s", string(sendData))

		rsp, err := dohttp.DoJsonHttp(map[string]string{}, "POST", c.callbackServer+router, sendData)
		if nil != err {
			logx.Errorf("post to callback error %s", err.Error())
			return err
		}
		defer rsp.Body.Close()
		return nil
	}
	logx.Infof("callback data is:%+v", data)
	return errors.New("callback url is null")
}
