package controller

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

// 反代
func Proxy(c *gin.Context) {
	alistUrl := config.Alist.URL

	target, err := url.Parse(alistUrl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "无法解析Alist站点地址",
			"error":   err.Error(),
		})
		return
	}

	// 获取请求路径
	path := c.Param("path")
	if path == "/" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "啥也没有，你让我你啥?",
		})
		return
	}

	// 创建反代
	proxy := httputil.NewSingleHostReverseProxy(target)

	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host

		if path != "" {
			req.URL.Path = singleJoiningSlash(target.Path, path)
		} else {
			req.URL.Path = target.Path
		}

		if target.RawQuery != "" {
			if req.URL.RawQuery != "" {
				req.URL.RawQuery = target.RawQuery + "&" + req.URL.RawQuery
			} else {
				req.URL.RawQuery = target.RawQuery
			}
		}

		if _, ok := req.Header["Host"]; !ok {
			req.Host = target.Host
		}

		req.Header.Set("X-Forwarded-Host", req.Host)
		req.Header.Set("X-Forwarded-Proto", req.URL.Scheme)
		req.Header.Set("X-Forwarded-For", c.ClientIP())
	}

	// 处理代理错误
	proxy.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, err error) {
		c.JSON(http.StatusBadGateway, gin.H{
			"code":    502,
			"message": "下载失败了，好可惜，看看错误提示吧",
			"error":   err.Error(),
		})
	}

	// 执行代理请求
	proxy.ServeHTTP(c.Writer, c.Request)
}

// 辅助函数：连接两个路径，避免双斜杠问题
func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")

	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
