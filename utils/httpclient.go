package utils

import (
	"bytes"
	"io"
	"net/http"
)

// HttpClient 发送http请求
func HttpClient(url string, method string, body []byte, header map[string]string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return "", err
	}

	for key, value := range header {
		req.Header.Set(key, value)
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
