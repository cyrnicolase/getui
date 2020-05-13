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

const (
	// KeyRegion 省市
	KeyRegion ConditionKey = `region`
	// KeyPhoneType 手机类型
	KeyPhoneType ConditionKey = `phonetype`
	// KeyTag 用户标签tag
	KeyTag ConditionKey = `tag`

	// OptTypeOr 取参数并集
	OptTypeOr OptType = 0
	// OptTypeAnd 取参数交集
	OptTypeAnd OptType = 1
	// OptTypeNotIn 相当于 not in
	OptTypeNotIn OptType = 2
)

// PushToApper 给App推送接口
type PushToApper interface {
	// PushToApp 给App推送消息
	PushToApp(PushToAppParam) ([]byte, error)

	// PushToAppContext 给App推送消息，并使用context控制
	PushToAppContext(PushToAppParam) ([]byte, error)
}

// ConditionKey 筛选条件
type ConditionKey string

// OptType 筛选组合类型
type OptType int

// Condition 筛选目标用户条件
type Condition struct {
	Key     ConditionKey `json:"key"`
	Values  []string     `json:"values"`
	OptType OptType      `json:"opt_type"`
}

// PushToAppParam 给应用推送消息参数
type PushToAppParam struct {
	RequestID     string        `json:"requestid"`
	Message       Message       `json:"message"`
	Link          *Link         `json:"link,omitempty"`
	Notification  *Notification `json:"notification,omitempty"`
	Transmission  *Transmission `json:"transmission,omitempty"`
	PushInfo      PushInfo      `json:"push_info,omitempty"`
	Condition     Condition     `json:"condition,omitempty"`
	Speed         int           `json:"speed,omitempty"`
	TaskName      string        `json:"task_name,omitempty"`
	DurationBegin string        `json:"duration_begin,omitempty"`
	DurationEnd   string        `json:"duration_end,omitempty"`
}

// PushToApp 给App推送消息
func (g *Getui) PushToApp(p PushToAppParam) ([]byte, error) {
	return g.PushToAppContext(context.Background(), p)
}

// PushToAppContext 给App推送消息，并使用context控制
func (g *Getui) PushToAppContext(ctx context.Context, p PushToAppParam) ([]byte, error) {
	if err := g.RefreshTokenContext(ctx); nil != err {
		return nil, err
	}

	if "" == p.RequestID {
		p.RequestID = strconv.FormatInt(time.Now().UnixNano(), 10)
	}
	url := fmt.Sprintf(`%s/%s/push_app`, APIServer, g.AppID)
	body, err := json.Marshal(p)
	if nil != err {
		return nil, errors.Wrap(err, "encode push_to_app param to json bytes")
	}

	return SendContext(ctx, url, g.Token, bytes.NewBuffer(body))
}
