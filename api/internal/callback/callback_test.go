package callback_test

import (
	"os"
	"testing"
	"wowfish/api/internal/callback"
)

func TestCallback(t *testing.T) {
	// data := &callback.CallBackToWowfishData{
	// 	CallbackBaseData: callback.CallbackBaseData{
	// 		Ret:  "0",
	// 		From: "123",
	// 		To:   "456",
	// 		Info: "",
	// 	},
	// 	Amount: "111",
	// }

	data := &callback.CallBackToWowfishNftData{
		CallbackBaseData: callback.CallbackBaseData{
			Ret:  0,
			From: "123",
			To:   "456",
			Info: "",
		},
		Id: "111",
	}

	err := callback.Instance().Callback("", data)
	t.Error(err)
}

func TestMain(m *testing.M) {
	callback.Instance().Init("111")
	exitCode := m.Run()

	os.Exit(exitCode)
}
