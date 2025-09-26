package main

import "gorm.io/gorm"

// 用户
type User struct {
	gorm.Model
	Username string `gorm:"unique"`
	Password string
}

// 帖子
type Post struct {
	gorm.Model
	UserID  uint
	Title   string
	Content string
}

// 评论
type Comment struct {
	gorm.Model
	PostID  uint
	UserID  uint
	Content string
}
