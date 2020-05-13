package getui

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

// Tokener 授权token请求接口
type Tokener interface {
	// RefreshToken 刷新认证Token
	RefreshToken() error

	// RefreshTokenContext 刷新token，并使用context进行控制
	RefreshTokenContext(context.Context) error
}

// RefreshToken 刷新授权Token
// 鉴权令牌有效期在个推默认是1天，过期后无法使用
// 所以在请求到新的Token后，会将该Token记录下来，并设置一个新的过期时间
func (g *Getui) RefreshToken() error {
	return g.RefreshTokenContext(context.Background())
}

// RefreshTokenContext 刷新Token，并使用context上下文进行限制
func (g *Getui) RefreshTokenContext(ctx context.Context) error {
	// 如果当前Token存在，且未过期，那么就不用刷新Token
	if "" != g.Token && time.Now().Before(g.TokenExpireAt) {
		return nil
	}

	now := time.Now()
	url := fmt.Sprintf(`%s/%s/auth_sign`, APIServer, g.AppID)
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
	output, err := SendContext(ctx, url, g.Token, bytes.NewBuffer(data))
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

func signature(appKey, timestamp, ms string) string {
	hash := sha256.New()
	hash.Write([]byte(appKey + timestamp + ms))
	sum := hash.Sum(nil)

	return fmt.Sprintf("%x", sum)
}
