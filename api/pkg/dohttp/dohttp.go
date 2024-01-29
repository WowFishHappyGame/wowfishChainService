package dohttp

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

// DoMultiFormHttp support multi-form body
func DoMultiFormHttp(headers map[string]string, Method, Url string, data map[string]string) (*http.Response, error) {
	var bufReader bytes.Buffer
	writer := multipart.NewWriter(&bufReader)

	for k, v := range data {
		_ = writer.WriteField(k, v)
	}

	if err := writer.Close(); err != nil {
		logx.Errorf("writer.Close:%s", err.Error())
		return nil, fmt.Errorf("writer.Close:%s", err.Error())
	}

	req, err := http.NewRequest(Method, Url, &bufReader)
	if err != nil {
		return nil, fmt.Errorf("%s %s NewRequest err: %s", Method, Url, err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	//set header
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	//logx.Infof("ContentLength=%d \n:%v", req.ContentLength, req)
	//GetMultiPart3(req)
	client := &http.Client{}
	rsp, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("http request %s err: %s", err, Url)
	}

	return rsp, nil
}

// DoTextHttp support text body
func DoTextHttp(headers map[string]string, Method, Url string, body string) (*http.Response, error) {
	var bufReader bytes.Buffer
	writer := io.StringWriter(&bufReader)

	_, err := writer.WriteString(body)
	if err != nil {
		return nil, fmt.Errorf("%s %s body write err: %s", Method, Url, err)
	}

	req, err := http.NewRequest(Method, Url, &bufReader)
	if err != nil {
		return nil, fmt.Errorf("%s %s NewRequest err: %s", Method, Url, err)
	}
	req.Header.Set("Content-Type", "text")

	//set header
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	client := &http.Client{}
	rsp, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("http request %s err: %s", err, Url)
	}

	return rsp, nil
}
