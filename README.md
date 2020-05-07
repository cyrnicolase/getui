# Go个推SDK

使用Go语言封装个推的RestAPI,自己提供一个SDK.

## 建议
将 getui 对象作为包级别变量. 因为Getui{} 本身会缓存请求接口的AuthToken，
在声明为全局变量的情况，缓存有效。

```
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