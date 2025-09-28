package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
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
	r.GET("/api/posts", func(c *gin.Context) {
		var posts []Post
		DB.Order("created_at desc").Find(&posts)
		c.JSON(200, posts)
	})

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

		// 帖子详情
		r.GET("/posts/:id", func(c *gin.Context) {
			id := c.Param("id")
			var post Post
			if err := DB.First(&post, id).Error; err != nil {
				c.JSON(404, gin.H{"error": "帖子未找到"})
				return
			}
			c.JSON(200, post)
		})

		// 获取帖子评论
		r.GET("/posts/:id/comments", func(c *gin.Context) {
			id := c.Param("id")
			var comments []Comment
			DB.Where("post_id = ?", id).Order("created_at desc").Find(&comments)
			c.JSON(200, comments)
		})
	}

	r.Run(":8082")
}
