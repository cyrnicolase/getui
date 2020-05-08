package getui

const (
	// MediaTypeImage 图片
	MediaTypeImage MediaType = iota + 1
	// MediaTypeAudio 音频
	MediaTypeAudio
	// MediaTypeVideo 视频
	MediaTypeVideo
)

// Template 模板
type Template interface {
	// Name 名字
	Name() string
}

// MediaType 媒体类型
type MediaType int

// Style 模板样式
// http://docs.getui.com/getui/server/rest/template/
type Style struct {
	Type         int    `json:"type"`
	Text         string `json:"text"`
	Title        string `json:"title"`
	BigStyle     int    `json:"big_style,omitempty"`
	BigImageURL  string `json:"big_image_url,omitempty"`
	Logo         string `json:"logo,omitempty"`
	LogoURL      string `json:"logourl,omitempty"`
	IsRing       bool   `json:"is_ring,omitempty"`
	IsVibrate    bool   `json:"is_vibrate,omitempty"`
	IsClearable  bool   `json:"is_clearable,omitempty"`
	NotifyID     int    `json:"notify_id,omitempty"`
	ChannelLevel int    `json:"channel_level,omitempty"`
}

// NewStyle 样式
// 默认: type=0
func NewStyle(text, title string) Style {
	return Style{
		Type:         0,
		Text:         text,
		Title:        title,
		Logo:         "",
		IsRing:       true,
		IsVibrate:    true,
		IsClearable:  true,
		BigStyle:     1,
		ChannelLevel: 4,
	}
}

// Notification 应用模板样式
// 模板消息的一种 (Template)
// http://docs.getui.com/getui/server/rest/template/
type Notification struct {
	TransmissionType    bool   `json:"transmission_type,omitempty"`
	TransmissionContent string `json:"transmission_content,omitempty"`
	DurationBegin       string `json:"duration_begin,omitempty"`
	DurationEnd         string `json:"duration_end,omitempty"`
	Style               Style  `json:"style"`
}

// Name 模板名字
func (Notification) Name() string {
	return "notification"
}

// NewNotification 返回消息类型
func NewNotification(text, title string) *Notification {
	return &Notification{
		TransmissionType: true,
		Style:            NewStyle(text, title),
	}
}

// Link 打开网页模板
// http://docs.getui.com/getui/server/rest/template/
type Link struct {
	URL           string `json:"url"`
	DurationBegin string `json:"duration_begin,omitempty"`
	DurationEnd   string `json:"duration_end,omitempty"`
	Style         Style  `json:"style"`
}

// Name 模板名字
func (Link) Name() string {
	return "link"
}

// NewLink 返回消息类型
func NewLink(url, text, title string) *Link {
	return &Link{
		URL:   url,
		Style: NewStyle(text, title),
	}
}

// Transmission 透传
// http://docs.getui.com/getui/server/rest/template/
type Transmission struct {
	TransmissionContent string `json:"transmission_content"`
	TransmissionType    bool   `json:"transmission_type,omitempty"`
	DurationBegin       string `json:"duration_begin,omitempty"`
	DurationEnd         string `json:"duration_end,omitempty"`
}

// Name 模板名字
func (Transmission) Name() string {
	return "transmission"
}

// NewTransmission 返回透传模板
// title 横幅标题
// text 横幅内容
// content 透传内容
func NewTransmission(text, title, content string) (*Transmission, PushInfo) {
	alert := Alert{}
	alert.Body = text
	alert.Title = title

	aps := APS{}
	aps.Alert = alert
	aps.AutoBadge = `+1`
	aps.ContentAvailable = 0

	pushInfo := PushInfo{}
	pushInfo.APS = aps
	pushInfo.Payload = content

	return &Transmission{
		TransmissionType:    true,
		TransmissionContent: content,
	}, pushInfo
}

// StartActivity 打开指定页面
// http://docs.getui.com/getui/server/rest/template/
type StartActivity struct {
	TransmissionType    bool   `json:"transmission_type,omitempty"`
	TransmissionContent string `json:"transmission_content,omitempty"`
	DurationBegin       string `json:"duration_begin,omitempty"`
	DurationEnd         string `json:"duration_end,omitempty"`
	Intent              string `json:"intent"`
	Style               Style  `json:"style"`
}

// Name 返回模板名字
func (StartActivity) Name() string {
	return "startactivity"
}

// NewStartActivity 返回打开指定页面的模板
func NewStartActivity(text, title string) *StartActivity {
	return &StartActivity{
		TransmissionType: true,
		Style:            NewStyle(text, title),
	}
}

// Alert 消息
type Alert struct {
	Title        string   `json:"title"`
	Body         string   `json:"body"`
	LocKey       string   `json:"loc-key,omitempty"`
	LocArgs      []string `json:"loc-args,omitempty"`
	LaunchImage  string   `json:"launch-image,omitempty"`
	TitleLocKey  string   `json:"title-loc-key,omitempty"`
	TitleLocArgs string   `json:"title-loc-args,omitempty"`
	ActionLocKey string   `json:"action-loc-key,omitempty"`
}

// APS 苹果推送
type APS struct {
	Alert            Alert  `json:"alert"`
	AutoBadge        string `json:"autoBadge"`
	ContentAvailable int    `json:"content-available"`
	Sound            string `json:"sound,omitempty"`
	Category         string `json:"category,omitempty"`
}

// Media 媒体信息
type Media struct {
	URL      string    `json:"url"`
	Type     MediaType `json:"type"`
	OnlyWifi bool      `json:"only_wifi,omitempty"`
}

// PushInfo IOS推送
type PushInfo struct {
	APS        APS     `json:"aps"`
	Payload    string  `json:"payload,omitempty"`
	Multimedia []Media `json:"multimedia,omitempty"`
}
