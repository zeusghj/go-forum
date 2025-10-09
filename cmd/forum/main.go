package main

import (
	"go-forum/internal/forum"
	"os"

	_ "go.uber.org/automaxprocs"
)

func main() {
	/*
		InitDB()
		r := gin.Default()

		// 解决跨域问题
		r.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"http://localhost:5173"},
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
			AllowCredentials: true,
		}))

		// 测试接口
		r.GET("/ping", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{"message": "pong"})
		})

		r.POST("/register", Register)
		r.POST("/login", login)

		// 列表帖子
		r.GET("/api/posts", ListPosts)

		// 需要鉴权的接口
		auth := r.Group("/api")
		auth.Use(AuthMiddleware())
		{
			auth.GET("/profile", func(ctx *gin.Context) {
				userID, _ := ctx.Get("user_id")
				username, _ := ctx.Get("username")
				ctx.JSON(http.StatusOK, gin.H{
					"user_id":  userID,
					"username": username,
				})
			})

			auth.POST("/posts", CreatePost)                  // 发帖接口
			auth.POST("/comments", CreateComment)            // 评论接口
			auth.GET("/posts/:id", GetPost)                  // 帖子详情
			auth.GET("/posts/:id/comments", GetPostComments) // 获取帖子评论
		}

		r.Run(":8082")
	*/

	command := forum.NewForumCommand()
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
