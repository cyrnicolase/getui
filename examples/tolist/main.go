package main

import (
	"encoding/json"
	"log"

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
	result, err := saveListBody()
	if nil != err {
		log.Fatal(err, "保存群发消息到服务器")
	}
	log.Println("保存群发消息到服务器", string(result))

	var body getui.ListBodyResponse
	if err := json.Unmarshal(result, &body); nil != err {
		log.Fatal(err, "decode bytes to SaveListBodyResponse")
	}
	if getui.ResultOK != body.Result {
		log.Fatal(err, "save_list_body api fail")
	}

	// 将保存到服务器的消息，群发给所有的客户端
	// 这里的群发可以使用 goroutine
	result, err = sendToList(body.TaskID)
	if nil != err {
		log.Fatal(err, "发送群发消息")
	}
	log.Println("发送群发消息", string(result))
}

func saveListBody() ([]byte, error) {
	msg := gt.NewMessage(getui.MsgTypeNotification)
	t := getui.NewNotification("群推消息1", "群推1")
	t.TransmissionContent = "群推透传消息"

	param := getui.ListBodyParam{
		Message:      msg,
		Notification: t,
	}

	return gt.SaveListBody(param)
}

func sendToList(taskid string) ([]byte, error) {
	p := getui.PushToListParam{
		Cid:        []string{Cid},
		TaskID:     taskid,
		NeedDetail: false,
	}

	return gt.PushToList(p)
}
