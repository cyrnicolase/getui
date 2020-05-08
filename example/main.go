package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/cyrnicolase/getui"
	"github.com/pkg/errors"
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
	if err := toSingle(); nil != err {
		log.Fatalf("%+v", err)
	}
}

func batchSingle() error {
	result, err := sendSingleBatch()
	if nil != err {
		return errors.Wrap(err, "批量推送单个消息")
	}
	log.Println("批量推送单个消息", string(result))

	return nil
}

func toList() error {
	result, err := saveListBody()
	if nil != err {
		return errors.Wrap(err, "保存群发消息到服务器")
	}
	log.Println("保存群发消息到服务器", string(result))

	var body getui.SaveListBodyResponse
	if err := json.Unmarshal(result, &body); nil != err {
		return errors.Wrap(err, "decode bytes to SaveListBodyResponse")
	}
	if getui.ResultOK != body.Result {
		return errors.Wrap(err, "save_list_body api fail")
	}

	// 将保存到服务器的消息，群发给所有的客户端
	// 这里的群发可以使用 goroutine
	result, err = sendToList(body.TaskID)
	if nil != err {
		return errors.Wrap(err, "发送群发消息")
	}
	log.Println("发送群发消息", string(result))

	return nil
}

func toSingle() error {
	result, err := sendTransmissionSingle()
	if nil != err {
		return errors.Wrap(err, "发送单个透传消息")
	}
	log.Println("发送单个透传消息", result)

	log.Println("==========")

	result, err = sendNotificationSingle()
	if nil != err {
		return errors.Wrap(err, "发送单个通知消息")
	}
	log.Println("发送单个通知消息", result)

	log.Println("==========")

	result, err = sendLinkSingle()
	if nil != err {
		return errors.Wrap(err, "发送单个连接消息")
	}
	log.Println("发送单个连接消息", result)

	return nil
}

func sendTransmissionSingle() (string, error) {
	msg := gt.NewMessage(getui.MsgTypeTransmission)
	t, pushInfo := getui.NewTransmission(`横幅内容`, `横幅标题`, "测试透传消息, time:"+time.Now().Format(`2006-01-02 15:04:05`))
	param := getui.PushToSingleParam{
		Message:      *msg,
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
	t, pushInfo := getui.NewTransmission("横幅标题", "横幅内容", "Batch测试透传消息, time:"+time.Now().Format(`15:04:05`))
	p := getui.NewPushToSingleParam(*msg)
	p.Transmission = t
	p.PushInfo = pushInfo
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
