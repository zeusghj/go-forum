package model

import "gorm.io/gorm"

type PostM struct {
	gorm.Model
	UserID  uint
	Title   string `gorm:"size:256;not null"`
	Content string `gorm:"not null"`
}

func (p *PostM) TableName() string {
	return "posts"
}
