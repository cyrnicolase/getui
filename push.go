package getui

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

func signature(appKey, timestamp, ms string) string {
	hash := sha256.New()
	hash.Write([]byte(appKey + timestamp + ms))
	sum := hash.Sum(nil)

	return fmt.Sprintf("%x", sum)
}

// RefreshToken 刷新授权Token
func (g *Getui) RefreshToken() error {
	// 如果当前Token存在，且未过期，那么就不用刷新Token
	if "" != g.Token && time.Now().Before(g.TokenExpireAt) {
		return nil
	}

	now := time.Now()
	url := fmt.Sprintf(`%s/%s/auth_sign`, RestAPI, g.AppID)
	timestamp := now.Unix() * 1000
	sign := signature(g.AppKey, strconv.FormatInt(timestamp, 10), g.MasterSecret)

	body := struct {
		Sign      string `json:"sign"`
		Timestamp int64  `json:"timestamp"`
		AppKey    string `json:"appkey"`
	}{
		Sign:      sign,
		Timestamp: timestamp,
		AppKey:    g.AppKey,
	}

	data, _ := json.Marshal(body)
	output, err := Send(url, g.Token, bytes.NewBuffer(data))
	if nil != err {
		return err
	}

	var resp struct {
		Result    string `json:"result"`
		AuthToken string `json:"auth_token"`
	}
	if err := json.Unmarshal(output, &resp); nil != err {
		return errors.Wrap(err, "parse json data to struct")
	}
	if ResultOK == resp.Result {
		g.Token = resp.AuthToken
		g.TokenExpireAt = now.Add(time.Hour * 23) // 设置为23个小时的有效期

		return nil
	}

	return errors.Errorf("refresh token. source result: %s", string(output))
}

// PushToSingleParam 单推
type PushToSingleParam struct {
	RequestID    string                 `json:"requestid"`
	Message      Message                `json:"message"`
	Link         *Link                  `json:"link,omitempty"`
	Notification *Notification          `json:"notification,omitempty"`
	Transmission *Transmission          `json:"transmission,omitempty"`
	PushInfo     map[string]interface{} `json:"push_info,omitempty"`
	Cid          string                 `json:"cid,omitempty"`
	Alias        string                 `json:"alias,omitempty"`
}

// NewPushToSingleParam 返回单推消息参数
func NewPushToSingleParam(message Message) *PushToSingleParam {
	return &PushToSingleParam{
		Message:   message,
		RequestID: strconv.FormatInt(time.Now().UnixNano(), 10),
	}
}

// PushToSingle 单推
func (g *Getui) PushToSingle(p PushToSingleParam) ([]byte, error) {
	// 刷新验证Token
	if err := g.RefreshToken(); nil != err {
		return nil, err
	}

	p.RequestID = strconv.FormatInt(time.Now().UnixNano(), 10)
	body, err := json.Marshal(p)
	if nil != err {
		return nil, errors.Wrap(err, "encode push_to_single param to json bytes")
	}
	url := fmt.Sprintf(`%s/%s/push_single`, RestAPI, g.AppID)

	return Send(url, g.Token, bytes.NewBuffer(body))
}

// SaveListBodyParam 保存消息体
type SaveListBodyParam struct {
	Message      Message                `json:"message"`
	Link         *Link                  `json:"link,omitempty"`
	Notification *Notification          `json:"notification,omitempty"`
	Transmission *Transmission          `json:"transmission,omitempty"`
	PushInfo     map[string]interface{} `json:"push_info,omitempty"`
	TaskName     string                 `json:"task_name,omitempty"`
}

// SaveListBodyResponse 保存群发消息服务器返回结构体
type SaveListBodyResponse struct {
	Result string `json:"result"`
	TaskID string `json:"taskid"`
}

// SaveListBody 将推送消息保存在服务器
// 后面可以重复调用tolist接口将保存的消息发送给不同的目标用户
func (g *Getui) SaveListBody(p SaveListBodyParam) ([]byte, error) {
	if err := g.RefreshToken(); nil != err {
		return nil, err
	}

	url := fmt.Sprintf(`%s/%s/save_list_body`, RestAPI, g.AppID)
	body, err := json.Marshal(p)
	if nil != err {
		return nil, errors.Wrap(err, "encode save_list_body param to json bytes")
	}

	return Send(url, g.Token, bytes.NewBuffer(body))
}

// PushListParam 群推参数
type PushListParam struct {
	Cid        []string `json:"cid"`
	Alias      []string `json:"alias"`
	TaskID     string   `json:"taskid"`
	NeedDetail bool     `json:"need_detail"`
}

// PushToList 群推
// 消息群发给cid list 或者 alias list列表对应的客户群,当两者并存的时候，以cid为准
// 并使用save_list_body返回的taskid,调用toList 接口,完成推送
func (g *Getui) PushToList(p PushListParam) ([]byte, error) {
	if err := g.RefreshToken(); nil != err {
		return nil, err
	}

	url := fmt.Sprintf(`%s/%s/push_list`, RestAPI, g.AppID)
	body, err := json.Marshal(p)
	if nil != err {
		return nil, errors.Wrap(err, "encode push_list param to json bytes")
	}

	return Send(url, g.Token, bytes.NewBuffer(body))
}

// PushSingleBatchParam 批量单推参数
// 将一批推送个人信息合并在一个接口，一次请求进行推送
type PushSingleBatchParam struct {
	MsgList    []PushToSingleParam `json:"msg_list"`
	NeedDetail bool                `json:"need_detail,omitempty"`
}

// PushSingleBatch 批量单推
func (g *Getui) PushSingleBatch(p PushSingleBatchParam) ([]byte, error) {
	if err := g.RefreshToken(); nil != err {
		return nil, err
	}

	url := fmt.Sprintf(`%s/%s/push_single_batch`, RestAPI, g.AppID)
	body, err := json.Marshal(p)
	if nil != err {
		return nil, errors.Wrap(err, "encode push_single_batch param to json bytes")
	}

	return Send(url, g.Token, bytes.NewBuffer(body))
}
