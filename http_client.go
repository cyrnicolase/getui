package getui

import (
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

// Send 请求
func Send(url, token string, body io.Reader) ([]byte, error) {
	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	req, err := http.NewRequest(`POST`, url, body)
	if nil != err {
		return nil, errors.Wrap(err, "generate http request")
	}

	// 处理请求新Token的请求
	if "" != token {
		req.Header.Add("authtoken", token)
	}
	req.Header.Add("Charset", "UTF-8")
	req.Header.Add("Content-Type", ContentTypeJSON)

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
