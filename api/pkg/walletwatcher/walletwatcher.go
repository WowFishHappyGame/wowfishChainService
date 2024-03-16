package walletwatcher

import (
	"context"
	"encoding/json"
	"io"
	"math"
	"strconv"
	"sync"
	"wowfish/api/internal/chain"
	"wowfish/api/internal/util"
	"wowfish/api/pkg/dohttp"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"
)

type WalletWatcher struct {
	emailServer string
	toEmail     string
	threshold   float64
	signKey     string
}

var kmsInstance *WalletWatcher
var once sync.Once

func Instance() *WalletWatcher {
	once.Do(func() {
		kmsInstance = &WalletWatcher{}
	})
	return kmsInstance
}

func (l *WalletWatcher) Init(emailServer string,
	toEmail string,
	threshold float64,
	signKey string) {
	l.emailServer = emailServer
	l.toEmail = toEmail
	l.threshold = threshold
	l.signKey = signKey
}

func (l *WalletWatcher) CheckAddressAmount(address string) {
	if l.emailServer != "" {
		chainMgr := chain.ChainClientInstance()
		balance, err := chainMgr.Provider.BalanceAt(context.Background(), common.HexToAddress(address), nil)
		if err != nil {
			logx.Error("query blance error ", err.Error())
			return
		}
		balanceInt := balance.Int64()
		if balanceInt < int64(l.threshold*math.Pow10(18)) {
			//send the warning email
			emailData := map[string]string{
				"email":   l.toEmail,
				"name":    "Wowfish",
				"subject": "Wallet Balance not enough",
				"body":    address + ":" + strconv.FormatFloat(float64(balanceInt)/math.Pow10(18), 'f', 10, 64),
			}

			emailDataToSign := map[string]any{}
			for key, value := range emailData {
				emailDataToSign[key] = value
			}
			encodeData := util.Instance().EncodeSignData(emailDataToSign)
			encodeData += "&" + l.signKey
			sign := util.Instance().ToMD5(encodeData)

			emailData["sign"] = sign

			rsp, err := dohttp.DoMultiFormHttp(map[string]string{}, "POST",
				l.emailServer,
				emailData)
			if err != nil {
				logx.Error("Send Metricsati email error", err.Error())
				return
			}
			defer rsp.Body.Close()

			body, err := io.ReadAll(rsp.Body)
			if err != nil {
				logx.Error("Send Metricsati email read response", err.Error())
				return
			}

			type EmailResponse struct {
				Code int64 `json:"code"`
			}

			var b = EmailResponse{}

			err = json.Unmarshal(body, &b)
			if err != nil {
				logx.Error("Send Metricsati Unmarshal response", err.Error())
				return
			}
			if b.Code != 1 {
				logx.Error("Send Metricsati email errcode", b.Code)
				return
			}
			logx.Info("send metrics email sucess")
		}
	}
}
