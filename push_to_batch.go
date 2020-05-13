package getui

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

// PushToBatcher 一次推送多个单推消息
type PushToBatcher interface {
	// PushToBatchSingle 将多个单推消息一次推送
	PushToBatchSingle(PushToBatchSingleParam) ([]byte, error)

	// PushToBatchSingleContext 批量单推，并携带上下文
	PushToBatchSingleContext(context.Context, PushToBatchSingleParam) ([]byte, error)
}

// PushToBatchSingleParam 批量单推参数
// 将一批推送个人信息合并在一个接口，一次请求进行推送
type PushToBatchSingleParam struct {
	MsgList    []PushToSingleParam `json:"msg_list"`
	NeedDetail bool                `json:"need_detail,omitempty"`
}

// PushToBatchSingle 批量单推
func (g *Getui) PushToBatchSingle(p PushToBatchSingleParam) ([]byte, error) {
	return g.PushToBatchSingleContext(context.Background(), p)
}

// PushToBatchSingleContext 批量单推，并携带上下文
func (g *Getui) PushToBatchSingleContext(ctx context.Context, p PushToBatchSingleParam) ([]byte, error) {
	if err := g.RefreshTokenContext(ctx); nil != err {
		return nil, err
	}

	url := fmt.Sprintf(`%s/%s/push_single_batch`, APIServer, g.AppID)
	body, err := json.Marshal(p)
	if nil != err {
		return nil, errors.Wrap(err, "encode push_single_batch param to json bytes")
	}

	return SendContext(ctx, url, g.Token, bytes.NewBuffer(body))
}
