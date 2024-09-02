package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Zkeai/go_template/common/conf"
	"github.com/Zkeai/go_template/common/logger"
	"github.com/Zkeai/go_template/common/redis"
	"github.com/gin-gonic/gin"
	redisv8 "github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strings"
	"time"
)

type CustomClaims struct {
	Wallet    string `json:"wallet"`
	SessionID string `json:"session_id"`
	jwt.RegisteredClaims
}

type ValidateTokenRes struct {
	Claims *CustomClaims
	Token  string
}

var SecretKey = []byte("muyu##coin..baby")

// GenerateToken 生成 JWT
func GenerateToken(wallet string, sessionID string) (string, error) {
	claims := CustomClaims{
		Wallet:    wallet,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Minute)), // 设置过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "yuka.1*.",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(SecretKey)
}

// ValidateToken 验证并解析 JWT
func ValidateToken(tokenString string) (*ValidateTokenRes, error) {
	token, _ := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return SecretKey, nil
	})

	//if err != nil {
	//	return &ValidateTokenRes{}, err
	//}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		//检查黑名单
		isBlacklisted, err := redis.GetClient().Get(redis.Ctx, tokenString).Result()
		if err == nil && isBlacklisted == "blacklisted" {
			return nil, fmt.Errorf("token is blacklisted")
		}

		return &ValidateTokenRes{
			Claims: claims,
			Token:  "",
		}, nil
	} else {

		if claims.ExpiresAt.Time.Before(time.Now()) {

			//判断sessionID 是否存在
			wallet := claims.Wallet
			result, err := redis.GetClient().Get(redis.Ctx, wallet).Result()

			if errors.Is(err, redisv8.Nil) {

				return nil, fmt.Errorf("token已过期")
			} else if err != nil {
				// 其他 Redis 错误
				return nil, err
			}

			// 将 JSON 字符串解析为 Go 的结构体
			var sessionData SessionData
			err = json.Unmarshal([]byte(result), &sessionData)
			if err != nil {
				logger.Error("Failed to unmarshal JSON data: %v", err)
				return nil, fmt.Errorf("failed to unmarshal")
			}
			sessionID := sessionData.SessionID
			if sessionID != claims.SessionID {
				return nil, fmt.Errorf("token错误")
			}
			newToken, err := GenerateToken(wallet, sessionID)
			userData := SessionData{
				Role:      sessionData.Role,
				Status:    sessionData.Status,
				Token:     newToken,
				SessionID: sessionID,
			}
			// 将数据转换为 JSON 字符串
			jsonData, _ := json.Marshal(userData)

			_, err = redis.GetClient().Set(redis.Ctx, wallet, jsonData, time.Minute*10).Result()

			return &ValidateTokenRes{
				Claims: claims,
				Token:  newToken,
			}, nil

		}

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
		c.Set("wallet", claims.Claims.Wallet)
		if claims.Token != "" {
			c.Header("Authorization", "Bearer "+claims.Token)
		}
		// 继续处理请求
		c.Next()
	}
}
