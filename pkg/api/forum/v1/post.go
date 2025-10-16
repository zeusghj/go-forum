package v1

import "time"

// PostInfo 指定了博客的详细信息.
type PostInfo struct {
	Username  string    `json:"username,omitempty"`
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// GetPostResponse 指定了 `GET /v1/post/detail` 接口的返回参数.
type GetPostResponse PostInfo

// CreatePostRequest 指定了 `GET /v1/post/publish`  接口的返回参数.
type CreatePostRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// ListPostRequest 指定了 `GET /v1/post/list` 接口的请求参数.
type ListPostRequest struct {
	Offset int `form:"offset" binding:"omitempty,min=0,max=100"`
	Limit  int `form:"limit" binding:"omitempty,min=0"`
}

// ListPostResponse 指定了 `GET /v1/post/list` 接口的返回参数.
type ListPostResponse struct {
	TotalCount int64       `json:"totalCount"`
	Posts      []*PostInfo `json:"posts"`
}
