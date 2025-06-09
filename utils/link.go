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
		Url      string      `json:"url"`
		Name     string      `json:"name"`
		Size     int64       `json:"size"`
		IsDir    bool        `json:"is_dir"`
		Modified string      `json:"modified"`
		Created  string      `json:"created"`
		Sign     string      `json:"sign"`
		Thumb    string      `json:"thumb"`
		Type     int         `json:"type"`
		Hashinfo string      `json:"hashinfo"`
		HashInfo interface{} `json:"hash_info"`
		RawUrl   string      `json:"raw_url"`
		Readme   string      `json:"readme"`
		Header   string      `json:"header"`
		Provider string      `json:"provider"`
		Related  interface{} `json:"related"`
	} `json:"data"`
}

// 获取链接地址
func GetFsLink(path string) (string, error) {
	header := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": GetToken(),
	}

	pass, err := AESDecrypt(conf.GlobalConfig.Alist.Password)
	if err != nil {
		return "", err
	}

	body := map[string]string{
		"path":     path,
		"password": pass,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	data, err := HttpClient(conf.GlobalConfig.Alist.APIURL+"/api/fs/get", "POST", jsonBody, header)
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

	return resp.Data.RawUrl, nil

}
