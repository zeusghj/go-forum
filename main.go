package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	InitDB()
	r := gin.Default()

	// 测试接口
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "pong"})
	})

	r.POST("/register", Register)
	r.POST("/login", login)

	// 需要鉴权的接口
	// auth := r.Group("/api", AuthMiddleware())
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

		auth.POST("/posts", CreatePost)       // 发帖接口
		auth.POST("/comments", CreateComment) // 评论接口
	}

	r.Run(":8082")
}
