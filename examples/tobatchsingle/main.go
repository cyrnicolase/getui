package main

import  (
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
	msg := gt.NewMessage(getui.MsgTypeTransmission)
	t, pushInfo := getui.NewTransmission("横幅标题", "横幅内容", "Batch测试透传消息, time:"+time.Now().Format(`15:04:05`))
	p := getui.NewPushToSingleParam(msg)
	p.Transmission = t
	p.PushInfo = pushInfo
	p.Cid = Cid

	msg1 := gt.NewMessage(getui.MsgTypeNotification)
	t1 := getui.NewNotification("Batch-Body", "Batch-Head")
	t1.TransmissionContent = "Batch我是一个通知"
	p1 := getui.NewPushToSingleParam(msg1)
	p1.Notification = t1
	p1.Cid = Cid

	msg2 := gt.NewMessage(getui.MsgTypeLink)
	t2 := getui.NewLink("www.baidu.com", "Batch访问百度", "Batch百度")
	p2 := getui.NewPushToSingleParam(msg2)
	p2.Link = t2
	p2.Cid = Cid

	var ml []getui.PushToSingleParam
	ml = append(ml, p)
	ml = append(ml, p1)
	ml = append(ml, p2)

	batch := getui.PushToBatchSingleParam{
		MsgList:    ml,
		NeedDetail: false,
	}

	result, err := gt.PushToBatchSingle(batch)
	if nil != err {
		log.Fatal(err, "批量单推")
	}

	log.Println(string(result))
}