package v1

import "time"

// LoginRequest 指定了 `POST /login` 接口的请求参数
type LoginRequest struct {
	Username string `json:"username" binding:"required,alphanum,min=1,max=255"`
	Password string `json:"password" binding:"required,min=6,max=18"`
}

// LoginResponse 指定了 `POST /login` 接口的返回参数
type LoginResponse struct {
	Token string `json:"token"`
}

// ChangePasswordRequest 指定了 `POST /v1/user/change-password` 接口的请求参数.
type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required,min=6,max=18"`
	NewPassword string `json:"newPassword" binding:"required,min=6,max=18"`
}

// CreateUserRequest 指定了 `POST /register` 接口的请求参数.
type CreateUserRequest struct {
	Username string `json:"username" binding:"required,alphanum,min=1,max=255"`
	Password string `json:"password" binding:"required,min=6,max=18"`
	Nickname string `json:"nickname" binding:"omitempty,max=30"`
	Email    string `json:"email" binding:"omitempty,email"`
	Phone    string `json:"phone" binding:"omitempty,len=11"`
}

// GetUserResponse 指定了 `GET /v1/user/user-info` 接口的返回参数.
type GetUserResponse UserInfo

// UserInfo 指定了用户的详细信息.
type UserInfo struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Nickname  string    `json:"nickname"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	PostCount int64     `json:"postCount"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
