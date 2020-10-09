package getui

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

const (
	// defaultRequestTimeout 默认请求超时时间
	defaultRequestTimeout = 3 * time.Second
)

var (
	// 请求超时时间
	requestTimeout time.Duration
)

// SetRequestTimeout 设置http请求超时时间
func SetRequestTimeout(d time.Duration) {
	requestTimeout = d
}

// Send 发送http请求，请求个推接口
func Send(url, token string, body io.Reader) ([]byte, error) {
	return SendContext(context.Background(), url, token, body)
}

// SendContext 携带context的发送请求
func SendContext(ctx context.Context, url, token string, body io.Reader) ([]byte, error) {
	client := newClient()
	req, err := newRequest(ctx, url, token, body)
	if nil != err {
		return nil, err
	}

	return request(client, req)
}

// 封装http client
func newClient() *http.Client {
	timeout := defaultRequestTimeout
	if requestTimeout.Nanoseconds() > 0 {
		timeout = requestTimeout
	}
	return &http.Client{
		Timeout: timeout,
	}
}

// 封装request header
func newRequest(ctx context.Context, url, token string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, `POST`, url, body)
	if nil != err {
		return nil, errors.Wrap(err, "generate http request")
	}

	// 处理请求新Token的请求
	if "" != token {
		req.Header.Add("authtoken", token)
	}
	req.Header.Add("Charset", "UTF-8")
	req.Header.Add("Content-Type", `application/json`)

	return req, nil
}

// 发起请求
func request(client *http.Client, req *http.Request) ([]byte, error) {
	resp, err := client.Do(req)
	if nil != err {
		return nil, errors.Wrap(err, "do http request")
	}

	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	if nil != err {
		return nil, errors.Wrap(err, "read response body")
	}

	return result, nil
}
