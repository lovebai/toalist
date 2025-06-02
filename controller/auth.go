package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/lovebai/toalist/utils"
)

var (
	jwtSecret = []byte("toalist-jwt-secret") // 也是直接写死了
	username  string
	password  string
)

// LoginRequest 登录请求结构
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse 登录响应结构
type LoginResponse struct {
	Token string `json:"token"`
}

// Claims JWT声明结构
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// 初始化登录信息
func InitLoginInfo() {
	username = config.Login.Username
	password = config.Login.Password
}

// 处理文件管理登录请求
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	// 验证用户名和密码
	if req.Username != username || utils.MD5Encrypt(req.Password) != password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 创建JWT声明 48小时过期 真难搞
	claims := Claims{
		Username: req.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(48 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	// 生成JWT令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成令牌失败"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{Token: tokenString})
}

// 认证中间件，懒得分离了.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供认证令牌"})
			c.Abort()
			return
		}

		// 从Bearer token中提取token
		tokenString := authHeader
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		}

		// 解析和验证token
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的认证令牌"})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("username", claims.Username)
		c.Next()
	}
}
