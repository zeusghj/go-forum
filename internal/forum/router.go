package forum

import (
	"go-forum/internal/forum/controller/v1/post"
	"go-forum/internal/forum/controller/v1/user"
	"go-forum/internal/forum/store"
	"go-forum/internal/pkg/core"
	"go-forum/internal/pkg/errno"
	"go-forum/internal/pkg/log"

	mw "go-forum/internal/pkg/middleware"

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
	pc := post.New(store.S)

	// 注册
	g.POST("/register", uc.Register)
	// 登录
	g.POST("/login", uc.Login)

	// 创建 v1 路由分组
	v1 := g.Group("/v1")
	{
		// 创建 user 路由分组
		userv1 := v1.Group("/user", mw.Authn())
		{
			userv1.GET("", uc.GetUser)                         // 获取指定用户名的用户详情
			userv1.GET("/profile", uc.GetUser)                 // 个人详情
			userv1.POST("/change-password", uc.ChangePassword) // 修改密码
		}
		// 创建 post 路由分组
		postv1 := v1.Group("/post", mw.Authn())
		{
			postv1.POST("/publish", pc.Create) // 发布动态
			postv1.GET("/list", pc.List)       // 帖子列表
			postv1.GET("/detail", pc.Get)      // 帖子详情

			// 创建 comment 路由分组
			comment := postv1.Group("/comment")
			{
				comment.POST("/add", pc.AddComment)  // 添加评论
				comment.GET("/list", pc.CommentList) // 评论列表
			}
		}

	}

	return nil
}
