package main

import (
	"encoding/json"
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
	result, err := sendSingleBatch()
	if nil != err {
		log.Fatal("批量推送单个消息", err)
	}
	log.Println("批量推送单个消息", string(result))
}

func toList() {
	result, err := saveListBody()
	if nil != err {
		log.Fatal("保存群发消息到服务器", err)
	}
	log.Println("保存群发消息到服务器", string(result))

	var body getui.SaveListBodyResponse
	if err := json.Unmarshal(result, &body); nil != err {
		log.Fatal("decode bytes to SaveListBodyResponse", err)
	}

	if getui.ResultOK != body.Result {
		log.Fatal("save_list_body api fail", string(result))
	}

	result, err = sendToList(body.TaskID)
	if nil != err {
		log.Fatal("发送群发消息", err)
	}
	log.Println("发送群发消息", string(result))

}

func toSingle() {
	result, err := sendTransmissionSingle()
	if nil != err {
		log.Fatal("发送单个透传消息", err)
	}
	log.Println("发送单个透传消息", result)

	log.Println("==========")

	result, err = sendNotificationSingle()
	if nil != err {
		log.Fatal("发送单个通知消息", err)
	}
	log.Println("发送单个通知消息", result)

	log.Println("==========")

	result, err = sendLinkSingle()
	if nil != err {
		log.Fatal("发送单个连接消息", err)
	}
	log.Println("发送单个连接消息", result)

}

func sendTransmissionSingle() (string, error) {
	msg := gt.NewMessage(getui.MsgTypeTransmission)
	t := getui.NewTransmission("测试透传消息, time:" + time.Now().Format(`2006-01-02 15:04:05`))
	param := getui.PushToSingleParam{
		Message:      *msg,
		Transmission: t,
		Cid:          Cid,
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
		Message:      *msg,
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
		Message: *msg,
		Link:    t,
		Cid:     Cid,
	}

	result, err := gt.PushToSingle(param)
	if nil != err {
		return "", err
	}

	return string(result), nil
}

func saveListBody() ([]byte, error) {
	msg := gt.NewMessage(getui.MsgTypeNotification)
	t := getui.NewNotification("群推消息1", "群推1")
	t.TransmissionContent = "群推透传消息"

	param := getui.SaveListBodyParam{
		Message:      *msg,
		Notification: t,
	}

	return gt.SaveListBody(param)
}

func sendToList(taskid string) ([]byte, error) {
	p := getui.PushListParam{
		Cid:        []string{Cid},
		TaskID:     taskid,
		NeedDetail: false,
	}

	return gt.PushToList(p)
}

func sendSingleBatch() ([]byte, error) {
	msg := gt.NewMessage(getui.MsgTypeTransmission)
	t := getui.NewTransmission("Batch测试透传消息, time:" + time.Now().Format(`15:04:05`))
	p := getui.NewPushToSingleParam(*msg)
	p.Transmission = t
	p.Cid = Cid

	msg1 := gt.NewMessage(getui.MsgTypeNotification)
	t1 := getui.NewNotification("Batch-Body", "Batch-Head")
	t1.TransmissionContent = "Batch我是一个通知"
	p1 := getui.NewPushToSingleParam(*msg1)
	p1.Notification = t1
	p1.Cid = Cid

	msg2 := gt.NewMessage(getui.MsgTypeLink)
	t2 := getui.NewLink("www.baidu.com", "Batch访问百度", "Batch百度")
	p2 := getui.NewPushToSingleParam(*msg2)
	p2.Link = t2
	p2.Cid = Cid

	var ml []getui.PushToSingleParam
	ml = append(ml, *p)
	ml = append(ml, *p1)
	ml = append(ml, *p2)

	batch := getui.PushSingleBatchParam{
		MsgList:    ml,
		NeedDetail: false,
	}

	return gt.PushSingleBatch(batch)
}
