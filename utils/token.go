package utils

import (
	"encoding/json"
	"errors"
	"log/slog"
	"time"

	"github.com/lovebai/toalist/conf"
)

type LoginBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Token string `json:"token"`
	} `json:"data"`
}

var (
	token           string
	tokenExpireTime time.Time
)

// 获取token
func loginAndGetToken() (string, error) {
	name := conf.GlobalConfig.Alist.Username
	pass, err := AESDecrypt(conf.GlobalConfig.Alist.Password)
	if err != nil {
		return "", err
	}

	body := LoginBody{
		Username: name,
		Password: pass,
	}

	header := map[string]string{
		"Content-Type": "application/json",
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	data, err := HttpClient(conf.GlobalConfig.Alist.APIURL+"/api/auth/login", "POST", jsonBody, header)
	if err != nil {
		return "", err
	}

	var resp Response
	err = json.Unmarshal([]byte(data), &resp)
	if err != nil {
		return "", err
	}

	if resp.Code != 200 {
		return "", errors.New(resp.Message)
	}
	slog.Info("成功获取 Alist token.")
	return resp.Data.Token, nil
}

// 检查token是否过期
func RefreshToken() {
	go func() {
		for {
			t, err := loginAndGetToken()
			// slog.Info("自动 Token 刷新已启动")
			if err != nil {
				slog.Error("Failed to refresh token", "error", err)
			} else {
				token = t
				tokenExpireTime = time.Now().Add(time.Hour * 48)
				// slog.Info("Token refreshed", "token", token)
			}
			time.Sleep(time.Hour * 47)
		}
	}()
}

func GetToken() string {
	if time.Now().After(tokenExpireTime) {
		t, err := loginAndGetToken()
		if err == nil {
			token = t
			tokenExpireTime = time.Now().Add(time.Hour * 48)
		}
	}
	return token
}
