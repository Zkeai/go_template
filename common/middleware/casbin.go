package middleware

import (
	"encoding/json"
	"github.com/Zkeai/go_template/common/logger"
	"github.com/Zkeai/go_template/common/redis"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SessionData struct {
	Role      int    `json:"role"`
	Status    int    `json:"status"`
	Token     string `json:"token"`
	SessionID string `json:"session_id"`
}

// CasbinMiddleware Casbin 中间件
func CasbinMiddleware(e *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		wallet := c.GetString("wallet")
		result, err := redis.GetClient().Get(c, wallet).Result()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "非法请求"})
			return
		}

		// 将 JSON 字符串解析为 Go 的结构体
		var sessionData SessionData
		err = json.Unmarshal([]byte(result), &sessionData)
		if err != nil {
			logger.Error("Failed to unmarshal JSON data: %v", err)
		}
		types := sessionData.Role
		// 获取当前用户的角色或身份标识
		sub := getRole(types)
		// 获取请求的 URL 和 Method
		obj := c.Request.URL.Path
		act := c.Request.Method

		// Casbin 权限检查
		ok, err := e.Enforce(sub, obj, act)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "验证失败, 没有权限"})
			return
		}

		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "无访问权限"})
			return
		}

		c.Next()
	}
}

func getRole(userType int) string {
	switch userType {
	case 0:
		return "user"
	case 1:
		return "merchant"
	case 2:
		return "admin"
	default:
		return ""
	}
}
