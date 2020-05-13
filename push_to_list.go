package getui

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

// SaveListBodyer 保存推送消息体接口
type SaveListBodyer interface {
	// SaveListBody 保存推送消息体
	SaveListBody(ListBodyParam) ([]byte, error)
}

// PushToLister 推送给一批客户端
type PushToLister interface {
	// PushToList 推送给一批用户
	PushToList(PushToListParam) ([]byte, error)

	// PushToListContext 推送给一批用户，同时使用上下文控制
	PushToListContext(context.Context, PushToListParam) ([]byte, error)
}

// ListBodyParam 保存消息体
type ListBodyParam struct {
	Message      Message       `json:"message"`
	Link         *Link         `json:"link,omitempty"`
	Notification *Notification `json:"notification,omitempty"`
	Transmission *Transmission `json:"transmission,omitempty"`
	PushInfo     PushInfo      `json:"push_info,omitempty"`
	TaskName     string        `json:"task_name,omitempty"`
}

// ListBodyResponse 保存群发消息服务器返回结构体
type ListBodyResponse struct {
	Result string `json:"result"`
	TaskID string `json:"taskid"`
}

// SaveListBody 将推送消息保存在服务器
// 后面可以重复调用tolist接口将保存的消息发送给不同的目标用户
func (g *Getui) SaveListBody(p ListBodyParam) ([]byte, error) {
	if err := g.RefreshToken(); nil != err {
		return nil, err
	}

	url := fmt.Sprintf(`%s/%s/save_list_body`, APIServer, g.AppID)
	body, err := json.Marshal(p)
	if nil != err {
		return nil, errors.Wrap(err, "encode save_list_body param to json bytes")
	}

	return Send(url, g.Token, bytes.NewBuffer(body))
}

// PushToListParam 群推参数
type PushToListParam struct {
	Cid        []string `json:"cid"`
	Alias      []string `json:"alias"`
	TaskID     string   `json:"taskid"`
	NeedDetail bool     `json:"need_detail"`
}

// PushToList 群推
// 消息群发给cid list 或者 alias list列表对应的客户群,当两者并存的时候，以cid为准
// 并使用save_list_body返回的taskid,调用toList 接口,完成推送
func (g *Getui) PushToList(p PushToListParam) ([]byte, error) {
	return g.PushToListContext(context.Background(), p)
}

// PushToListContext 携带上下文的发送消息给一批用户
func (g *Getui) PushToListContext(ctx context.Context, p PushToListParam) ([]byte, error) {
	if err := g.RefreshTokenContext(ctx); nil != err {
		return nil, err
	}

	url := fmt.Sprintf(`%s/%s/push_list`, APIServer, g.AppID)
	body, err := json.Marshal(p)
	if nil != err {
		return nil, errors.Wrap(err, "encode push_list param to json bytes")
	}

	return SendContext(ctx, url, g.Token, bytes.NewBuffer(body))
}
