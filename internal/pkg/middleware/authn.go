package middleware

import (
	"go-forum/internal/pkg/core"
	"go-forum/internal/pkg/errno"
	"go-forum/internal/pkg/known"
	"go-forum/pkg/token"

	"github.com/gin-gonic/gin"
)

// Authn 是认证中间件，用来从 gin.Context 中提取 token 并验证 token 是否合法，
// 如果合法则将 token 中的 sub 作为<用户名>存放在 gin.Context 的 XUsernameKey 键中.
func Authn() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 解析 JWT Token 先获取token再从token取出username
		userID, username, err := token.ParseRequest(ctx)

		if err != nil || userID == 0 {
			core.WriteResponse(ctx, errno.ErrTokenInvalid, nil)
			ctx.Abort()

			return
		}

		ctx.Set(known.XUserIDKey, userID)
		ctx.Set(known.XUsernameKey, username)
		ctx.Next()
	}
}
