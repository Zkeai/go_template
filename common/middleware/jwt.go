package middleware

import (
	"fmt"
	"github.com/Zkeai/go_template/common/conf"
	"github.com/Zkeai/go_template/common/redis"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strings"
	"time"
)

type CustomClaims struct {
	Wallet string `json:"wallet"`
	jwt.RegisteredClaims
}

var SecretKey = []byte("muyu##coin..baby")

// GenerateToken 生成 JWT
func GenerateToken(wallet string) (string, error) {
	claims := CustomClaims{
		Wallet: wallet,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), // 设置过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "yuka.1*.",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(SecretKey)
}

// ValidateToken 验证并解析 JWT
func ValidateToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return SecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		//检查黑名单
		isBlacklisted, err := redis.GetClient().Get(redis.Ctx, tokenString).Result()
		if err == nil && isBlacklisted == "blacklisted" {
			return nil, fmt.Errorf("token is blacklisted")
		}

		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}

// InvalidateToken 将 JWT 加入黑名单
func InvalidateToken(tokenString string) error {
	return redis.GetClient().Set(redis.Ctx, tokenString, "blacklisted", time.Until(time.Now().Add(time.Hour*24))).Err()
}

// Middleware 是一个 Gin 中间件，用于验证 JWT
func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取 token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, conf.Response{Code: http.StatusUnauthorized, Msg: "Validate error", Data: gin.H{"error": "请求头中缺少 Authorization"}})
			c.Abort()
			return
		}

		// 检查 token 的前缀
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, conf.Response{Code: http.StatusUnauthorized, Msg: "Validate error", Data: gin.H{"error": "请求头中的 Authorization 格式错误"}})
			c.Abort()
			return
		}

		// 验证并解析 token
		claims, err := ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, conf.Response{Code: http.StatusUnauthorized, Msg: "Validate error", Data: gin.H{"error": err.Error()}})
			c.Abort()
			return
		}

		// 将解析后的用户 wallet 设置到上下文中
		c.Set("wallet", claims.Wallet)

		// 继续处理请求
		c.Next()
	}
}
