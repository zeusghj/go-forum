package model

import "gorm.io/gorm"

type CommentM struct {
	gorm.Model
	PostID  uint
	UserID  uint
	Content string `gorm:"not null"`
}

func (c *CommentM) TableName() string {
	return "comments"
}
