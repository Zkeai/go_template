package cttp

import (
	"context"
	"github.com/Zkeai/go_template/common/conf"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const userIDHeader = "Authorization"

type headerKey struct{}

// Header head信息.
type Header struct {
	JwtToken string `json:"token"`
}

// HeaderHandler header信息写入context.
func HeaderHandler(ginCtx *gin.Context) {
	h := &Header{}

	authHeader := ginCtx.GetHeader(userIDHeader)
	if authHeader == "" {
		ginCtx.JSON(http.StatusUnauthorized, conf.Response{Code: http.StatusUnauthorized, Msg: "Validate error", Data: gin.H{"error": "请求头中缺少 Authorization"}})
		ginCtx.Abort()
		return
	}
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	h.JwtToken = tokenStr

	ctx := ginCtx.Request.Context()

	ctx = context.WithValue(ctx, headerKey{}, h)

	ginCtx.Request = ginCtx.Request.WithContext(ctx)
	ginCtx.Next()
}

// GetHeader 获取head信息.
func GetHeader(ctx context.Context) *Header {
	h, _ := ctx.Value(headerKey{}).(*Header)
	return h
}

// GetJwtToken 获取Uuid.
func GetJwtToken(ctx context.Context) string {
	h := GetHeader(ctx)
	if h == nil {
		return ""
	}

	return h.JwtToken
}
