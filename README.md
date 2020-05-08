# Go个推SDK

使用Go语言封装个推的RestAPI,自己提供一个SDK.

## 支持

+ 对使用App的某个用户，单独推送消息(push_single)
+ tolist群推(push_list)
+ 批量单推(push_single_batch)

## 用法

```go
gt := getui.NewGetui(
    os.Getenv("APP_ID"),
    os.Getenv("APP_SECRET"),
    os.Getenv("APP_KEY"),
    os.Getenv("MASTER_SECRET"),
)

cid := `clientid`
message := gt.NewMessage(getui.MsgTypeTransmission)
template, pushInfo := getui.NewTransmission(`横幅内容`, `横幅标题`, `透传内容`)
param := getui.PushToSingleParam {
    Message:        *message,
    Transmission:   template,
    Cid:            cid,
    PushInfo:       pushInfo,
}

result, err := gt.PushToSingle(param)
if nil != err {
    // handle error
}

fmt.Println(string(result))
```

其他使用方法请查看[example](https://github.com/cyrnicolase/getui/blob/master/example/main.go)

## 建议
将 getui 对象作为包级别变量. 因为Getui{} 本身会缓存请求接口的AuthToken，
在声明为全局变量的情况，缓存有效。

```go
var (
    gt = &getui.Getui{}
)

func init() {
    gt = getui.NewGetui(
        os.Getenv("APP_ID"),
        os.Getenv("APP_SECRET"),
        os.Getenv("APP_KEY"),
        os.Getenv("MASTER_SECRET"),
    )
}

```