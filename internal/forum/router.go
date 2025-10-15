package forum

import (
	"go-forum/internal/forum/controller/v1/user"
	"go-forum/internal/forum/store"
	"go-forum/internal/pkg/core"
	"go-forum/internal/pkg/errno"
	"go-forum/internal/pkg/log"

	"github.com/gin-gonic/gin"
)

// installRouters 安装 miniblog 接口路由.
func installRouters(g *gin.Engine) error {
	// 注册 404 Handler
	g.NoRoute(func(ctx *gin.Context) {
		core.WriteResponse(ctx, errno.ErrPageNotFound, nil)
	})

	// 注册 /healthz handler.
	g.GET("/healthz", func(c *gin.Context) {
		log.C(c).Infow("Healthz function called")

		core.WriteResponse(c, nil, map[string]string{"status": "ok"})
	})

	uc := user.New(store.S)

	// 注册
	g.POST("/register", uc.Register)
	// 登录
	g.POST("/login", uc.Login)

	// 创建 v1 路由分组
	// v1 := g.Group("/v1")
	{

		// 创建 users 路由分组
		// userv1 := v1.Group("/user")
		// {
		// 	userv1.POST("", uc.Create)
		// }
	}

	return nil
}
