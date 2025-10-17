package model

import "gorm.io/gorm"

type PostM struct {
	gorm.Model
	UserID  uint
	Title   string `gorm:"size:256;not null"`
	Content string `gorm:"not null"`
}

// CommentCountResult 每个帖子的评论数量
type CommentCountResult struct {
	PostID uint  `gorm:"column:post_id"`
	Count  int64 `gorm:"column:count"`
}

func (p *PostM) TableName() string {
	return "posts"
}
