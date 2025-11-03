package middleware

import (
	"fmt"
	"go-forum/internal/pkg/core"
	"go-forum/internal/pkg/errno"
	"go-forum/internal/pkg/log"
	"go-forum/internal/pkg/model"

	"github.com/gin-gonic/gin"
)

type Auther interface {
	Authorize(sub, obj, act string) (bool, error)
}

// Authz 是 Gin 中间件，用来进行请求授权.
func Authz(a Auther) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// sub := ctx.GetString(known.XUsernameKey)
		v, exists := ctx.Get("user")
		if !exists {
			core.WriteResponse(ctx, errno.ErrTokenInvalid, nil)
			return
		}
		user := v.(model.UserM)
		obj := ctx.FullPath()     // 路径
		act := ctx.Request.Method // 动作

		fmt.Print("sub = ", user.Role)

		log.Debugw("Build authorize context", "sub", user.Username, "obj", obj, "act", act)
		if allowed, _ := a.Authorize(user.Username, obj, act); !allowed {
			core.WriteResponse(ctx, errno.ErrUnauthorized, nil)
			ctx.Abort()
			return
		}
	}
}
