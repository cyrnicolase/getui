package main

import (
	"log"
	"time"

	"github.com/cyrnicolase/getui"
)

var (
	// ID AppID
	ID = `ay187m5uhy6CqjdAbKh8dA`
	// Key AppKey
	Key = `sHmi1wyScd8Cvoq2hgQ1h1`
	// Secret AppSecret
	Secret = `fY5UUx8r8GA71u1tdZSRe6`
	// MasterSecret 主密钥
	MasterSecret = `4PaHfGXc4b6aPD9VBOWla1`
	// Cid 测试客户端clientid
	Cid = `6ee9804698cff40c159f59be12bd4bf1`
	// 声明为包级别变量，这样可以缓存接口请求Token
	gt = &getui.Getui{}
)

func init() {
	gt = getui.NewGetui(ID, Secret, Key, MasterSecret)
}

func main() {
	result, err := sendTransmissionSingle()
	if nil != err {
		log.Fatal(err, "发送单个透传消息")
	}
	log.Println("发送单个透传消息", result)

	log.Println("==========")

	result, err = sendNotificationSingle()
	if nil != err {
		log.Fatal(err, "发送单个通知消息")
	}
	log.Println("发送单个通知消息", result)

	log.Println("==========")

	result, err = sendLinkSingle()
	if nil != err {
		log.Fatal(err, "发送单个连接消息")
	}
	log.Println("发送单个连接消息", result)
}

func sendTransmissionSingle() (string, error) {
	msg := gt.NewMessage(getui.MsgTypeTransmission)
	t, pushInfo := getui.NewTransmission(`横幅内容`, `横幅标题`, "测试透传消息, time:"+time.Now().Format(`2006-01-02 15:04:05`))
	param := getui.PushToSingleParam{
		Message:      msg,
		Transmission: t,
		Cid:          Cid,
		PushInfo:     pushInfo,
	}

	result, err := gt.PushToSingle(param)
	if nil != err {
		return "", err
	}

	return string(result), nil
}

func sendNotificationSingle() (string, error) {
	msg := gt.NewMessage(getui.MsgTypeNotification)
	t := getui.NewNotification("Body", "Head")
	t.TransmissionContent = "我是一个通知"

	param := getui.PushToSingleParam{
		Message:      msg,
		Notification: t,
		Cid:          Cid,
	}

	result, err := gt.PushToSingle(param)
	if nil != err {
		return "", err
	}

	return string(result), nil
}

func sendLinkSingle() (string, error) {
	msg := gt.NewMessage(getui.MsgTypeLink)
	t := getui.NewLink("www.baidu.com", "访问百度", "百度")

	param := getui.PushToSingleParam{
		Message: msg,
		Link:    t,
		Cid:     Cid,
	}

	result, err := gt.PushToSingle(param)
	if nil != err {
		return "", err
	}

	return string(result), nil
}
