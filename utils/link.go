package utils

import (
	"encoding/json"
	"errors"

	"github.com/lovebai/toalist/conf"
)

type ResBody struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Url string `json:"url"`
	} `json:"data"`
}

// 获取链接地址
func GetFsLink(path string) (string, error) {
	header := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": GetToken(),
	}

	body := map[string]string{
		"path": path,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	data, err := HttpClient(conf.GlobalConfig.Alist.APIURL+"/api/fs/link", "POST", jsonBody, header)
	if err != nil {
		return "", err
	}

	var resp ResBody
	err = json.Unmarshal([]byte(data), &resp)
	if err != nil {
		return "", err
	}

	if resp.Code != 200 {
		return "", errors.New(resp.Message)
	}

	return resp.Data.Url, nil

}
