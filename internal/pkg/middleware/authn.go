package middleware

import (
	"fastgo/internal/pkg/contextx"
	"fastgo/internal/pkg/core"
	"fastgo/internal/pkg/errorsx"
	"fastgo/pkg/token"
	"github.com/gin-gonic/gin"
)

// Authn 为认证中间件, 该函数将从 gin.Context 中提取 token 并验证是否合法.
// 若 token 合法, 则从中解析出 userID 并将其注入上下文
func Authn() gin.HandlerFunc {
	return func(context *gin.Context) {
		// 解析 JWT Token
		userID, err := token.ParseRequest(context)
		if err != nil {
			core.WriteResponse(context, errorsx.ErrTokenInvalid, nil)
			context.Abort()
			return
		}

		// 解析成功, 将用户ID和用户名注入上下文
		ctx := contextx.WithUserID(context.Request.Context(), userID)
		context.Request = context.Request.WithContext(ctx)

		// 继续执行主线程
		context.Next()
	}
}
