package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"kowhai/global"
	"net/http"
	"strings"
	"time"
)

type Claims struct {
	Name string `json:"name"`
	jwt.StandardClaims
}

// jwt中间件
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求头中的token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		// 解析并验证jwt
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			// 校验jwt的签名算法
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				global.Logger.Error("Unexpected signing method: %v", token.Header["alg"])
				return nil, jwt.ErrSignatureInvalid
			}
			return mySigningKey, nil
		})
		// 如果解析失败，返回401
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// 如果token有效，将用户信息存放到Context
		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			c.Set("name", claims.Name)
		}
		c.Next()
	}
}

// 创建jwt token
func CreateToken(name string) (string, error) {
	claims := Claims{
		Name: name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(), //过期时间设置为3天
		},
	}
	//生成token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//使用密钥签名
	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		global.Logger.Error("generate token failed, error: %v", err)
		return "", err
	}
	return tokenString, nil
}
