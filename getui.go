package getui

import "time"

const (
	// APIServer 接口域名
	APIServer string = `https://restapi.getui.com/v1`

	// ResultOK 接口返回状态；成功状态
	ResultOK string = `ok`

	// MsgTypeNotification 打开应用模板
	MsgTypeNotification MsgType = `notification`
	// MsgTypeLink 打开通知网页模板
	MsgTypeLink MsgType = `link`
	// MsgTypeNotyPopload 通知弹窗下载模板
	MsgTypeNotyPopload MsgType = `notypopload`
	// MsgTypeStartActivity 打开应用内特定页面模板
	MsgTypeStartActivity MsgType = `startactivity`
	// MsgTypeTransmission 透传消息模板
	MsgTypeTransmission MsgType = `transmission`

	// NetWorkTypeUnlimited 联网方式不限
	NetWorkTypeUnlimited NetWorkType = 0
	// NetWorkTypeWIFI 仅限wifi
	NetWorkTypeWIFI NetWorkType = 1
	// NetWorkTypeMobile 仅移动网
	NetWorkTypeMobile NetWorkType = 2
)

// MsgType 接收消息类型
type MsgType string

// NetWorkType 网络形式
type NetWorkType int

// Getui 应用配置
// 这个建议调用方做为包级别变量，这样可以进行数据缓存，缓存Token 信息
// http://docs.getui.com/getui/server/rest/explain/
type Getui struct {
	// AppID 个推分配应用ID
	AppID string

	// AppSecret 个推分配应用加密
	AppSecret string

	// AppKey 个推分配Key
	AppKey string

	// MasterSecret 个推分配主加密串
	MasterSecret string

	// Token 请求个推需要的授权Token
	Token string

	// TokenExpireAt 过期时间点
	// 这个是用来记录Token的过期时间，如果超过这个时间点
	// 那么就需要重新找个推请求新的Token
	TokenExpireAt time.Time
}

// NewGetui 返回个推结构体；并使用该对象来调用相关推送方法
func NewGetui(appID, appSecret, appKey, masterSecret string) *Getui {
	return &Getui{
		AppID:        appID,
		AppSecret:    appSecret,
		AppKey:       appKey,
		MasterSecret: masterSecret,
	}
}

// Message 消息体类型
// http://docs.getui.com/getui/server/rest/explain/
type Message struct {
	AppKey            string      `json:"appkey"`
	IsOffline         bool        `json:"is_offline"`
	OfflineExpireTime int         `json:"offline_expire_time"`
	PushNetworkType   NetWorkType `json:"push_network_type"`
	MsgType           MsgType     `json:"msgtype"`
}

// NewMessage 返回消息类型
// 返回消息类型；
// 默认：is_offline = true; offline_expire_time=86400000; push_network_type = 0
func (g Getui) NewMessage(msgType MsgType) Message {
	return Message{
		AppKey:            g.AppKey,
		IsOffline:         true,
		OfflineExpireTime: 86400000,
		PushNetworkType:   NetWorkTypeUnlimited,
		MsgType:           msgType,
	}
}
