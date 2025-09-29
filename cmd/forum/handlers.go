package main

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

var jwtSecret = []byte("Rtg8BPKNEf2mB4mgvKONGPZZQSaJWNLijxR42qRgq0iBb5")

// 注册接口
func Register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	user := User{
		Username: req.Username,
		Password: req.Password,
	}

	// 插入数据库
	if err := DB.Create(&user).Error; err != nil {
		// 可能是用户名重复
		if err == gorm.ErrDuplicatedKey {
			c.JSON(http.StatusConflict, gin.H{"error": "用户名已存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "注册成功", "user_id": user.ID})
}

// 登录接口
func login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var user User
	if err := DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户不存在"})
		return
	}

	if user.Password != req.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "密码错误"})
		return
	}

	// 生成 JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // 一天过期
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成 token 失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "登录成功", "token": tokenString})
}

// 发帖接口
func CreatePost(c *gin.Context) {
	// 从上下文获取 user_id （JWT中间件存的）
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未登录"})
		return
	}
	// userID := userIDValue.(uint)

	userIDFloat, ok := userIDValue.(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户ID解析失败"})
		return
	}
	userID := uint(userIDFloat)

	// 解析请求
	var req struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	post := Post{
		UserID:  userID,
		Title:   req.Title,
		Content: req.Content,
	}

	if err := DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建帖子失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "发帖成功",
		"post_id": post.ID,
	})
}

// 帖子列表 GET /api/posts
func ListPosts(c *gin.Context) {
	var posts []Post
	if err := DB.Order("created_at desc").Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询帖子失败"})
		return
	}

	// 收集 post ids 和 user ids
	postIDs := make([]uint, 0, len(posts))
	userIDSet := make(map[uint]struct{})
	userIDs := make([]uint, 0, len(posts))
	for _, p := range posts {
		postIDs = append(postIDs, p.ID)
		if _, ok := userIDSet[p.UserID]; !ok {
			userIDSet[p.UserID] = struct{}{}
			userIDs = append(userIDs, p.UserID)
		}
	}

	// 查询用户信息 （批量）
	userMap := map[uint]string{}
	if len(userIDs) > 0 {
		var users []User
		if err := DB.Where("id IN ?", userIDs).Find(&users).Error; err == nil {
			for _, u := range users {
				userMap[u.ID] = u.Username
			}
		}
	}

	// 查询每个帖子的评论数量(批量)
	type countResult struct {
		PostID uint  `gorm:"column:post_id"`
		Count  int64 `gorm:"column:count"`
	}
	commentCountMap := map[uint]int64{}
	if len(postIDs) > 0 {
		var counts []countResult
		DB.Model(&Comment{}).
			Select("post_id, count(*) as count").
			Where("post_id IN ?", postIDs).
			Group("post_id").
			Find(&counts)
		for _, ct := range counts {
			commentCountMap[ct.PostID] = ct.Count
		}
	}

	// 构造响应
	type PostResp struct {
		ID           uint      `json:"id"`
		Title        string    `json:"title"`
		Content      string    `json:"content"`
		UserID       uint      `json:"user_id"`
		Username     string    `json:"username"`
		CreatedAt    time.Time `json:"created_at"`
		CommentCount int64     `json:"comment_count"`
	}

	resp := make([]PostResp, 0, len(posts))
	for _, p := range posts {
		resp = append(resp, PostResp{
			ID:           p.ID,
			Title:        p.Title,
			Content:      p.Content,
			UserID:       p.UserID,
			Username:     userMap[p.UserID], // 若为空，前端可处理显示“匿名”或 id
			CreatedAt:    p.CreatedAt,
			CommentCount: commentCountMap[p.ID],
		})
	}

	c.JSON(http.StatusOK, resp)

}

// 帖子详情 GET /api/posts/:id
func GetPost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "帖子ID不合法"})
		return
	}

	var post Post
	if err := DB.First(&post, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "帖子未找到"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询帖子失败"})
		return
	}

	// 查询作者用户名(可选)
	var author User
	authorName := ""
	if err := DB.First(&author, post.UserID).Error; err == nil {
		authorName = author.Username
	}

	// 查询评论数
	var commentCount int64
	DB.Model(&Comment{}).Where("post_id = ?", post.ID).Count(&commentCount)

	resp := gin.H{
		"id":            post.ID,
		"title":         post.Title,
		"content":       post.Content,
		"user_id":       post.UserID,
		"username":      authorName,
		"created_at":    post.CreatedAt,
		"updated_at":    post.UpdatedAt,
		"comment_count": commentCount,
	}

	c.JSON(http.StatusOK, resp)
}

// 评论帖子
func CreateComment(c *gin.Context) {
	// 1. 获取用户 ID
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}
	userIDFloat, ok := userIDValue.(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户ID解析失败"})
		return
	}
	userID := uint(userIDFloat)

	// 2. 获取请求体
	var req struct {
		PostID  uint   `json:"post_id" binding:"required"`
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3. 新建评论
	comment := Comment{
		PostID:  req.PostID,
		UserID:  userID,
		Content: req.Content,
	}

	if err := DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "评论失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "评论成功",
		"comment": comment,
	})

}

// 获取帖子评论 GET /api/posts/:id/comments
func GetPostComments(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "帖子ID不合法"})
		return
	}

	var comments []Comment
	if err := DB.Where("post_id = ?", id).Order("created_at asc").Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询评论失败"})
		return
	}

	// 批量查询评论用户的用户名
	userIDSet := map[uint]struct{}{}
	userIDs := make([]uint, 0, len(comments))
	for _, cm := range comments {
		if _, ok := userIDSet[cm.UserID]; !ok {
			userIDSet[cm.UserID] = struct{}{}
			userIDs = append(userIDs, cm.UserID)
		}
	}

	userMap := map[uint]string{}
	if len(userIDs) > 0 {
		var users []User
		if err := DB.Where("id IN ?", userIDs).Find(&users).Error; err == nil {
			for _, u := range users {
				userMap[u.ID] = u.Username
			}
		}
	}

	type CommentResp struct {
		ID        uint      `json:"id"`
		PostID    uint      `json:"post_id"`
		UserID    uint      `json:"user_id"`
		Username  string    `json:"username"`
		Content   string    `json:"content"`
		CreatedAt time.Time `json:"created_at"`
	}

	res := make([]CommentResp, 0, len(comments))
	for _, cm := range comments {
		res = append(res, CommentResp{
			ID:        cm.ID,
			PostID:    cm.PostID,
			UserID:    cm.UserID,
			Username:  userMap[cm.UserID],
			Content:   cm.Content,
			CreatedAt: cm.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, res)
}
