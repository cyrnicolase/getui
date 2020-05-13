package getui

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

// PushToSingler 推送单个消息到指定客户端
type PushToSingler interface {
	// 推送单个消息到指定客户端
	PushToSingle(PushToSingleParam) ([]byte, error)

	// PushToSingleContext 携带上下文的推送单个消息到指定客户端
	PushToSingleContext(context.Context, PushToSingleParam) ([]byte, error)
}

// PushToSingleParam 单推
type PushToSingleParam struct {
	RequestID    string        `json:"requestid"`
	Message      Message       `json:"message"`
	Link         *Link         `json:"link,omitempty"`
	Notification *Notification `json:"notification,omitempty"`
	Transmission *Transmission `json:"transmission,omitempty"`
	PushInfo     PushInfo      `json:"push_info,omitempty"`
	Cid          string        `json:"cid,omitempty"`
	Alias        string        `json:"alias,omitempty"`
}

// NewPushToSingleParam 返回单推消息参数
func NewPushToSingleParam(message Message) PushToSingleParam {
	return PushToSingleParam{
		Message:   message,
		RequestID: strconv.FormatInt(time.Now().UnixNano(), 10),
	}
}

// PushToSingle 单推
func (g *Getui) PushToSingle(p PushToSingleParam) ([]byte, error) {
	return g.PushToSingleContext(context.Background(), p)
}

// PushToSingleContext 带有上下文的推送
func (g *Getui) PushToSingleContext(ctx context.Context, p PushToSingleParam) ([]byte, error) {
	if err := g.RefreshTokenContext(ctx); nil != err {
		return nil, err
	}

	if "" == p.RequestID {
		p.RequestID = strconv.FormatInt(time.Now().UnixNano(), 10)
	}
	url := fmt.Sprintf(`%s/%s/push_single`, APIServer, g.AppID)
	body, err := json.Marshal(p)
	if nil != err {
		return nil, errors.Wrap(err, "encode push_to_single param to json bytes")
	}

	return SendContext(ctx, url, g.Token, bytes.NewBuffer(body))
}
