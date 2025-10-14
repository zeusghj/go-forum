package model

import "gorm.io/gorm"

type UserM struct {
	gorm.Model
	Username string `gorm:"unique;size:255;not null"`
	Password string `gorm:"size:255;not null"`
	Nickname string `gorm:"size:30"`
	Email    string `gorm:"size:256"`
	Phone    string `gorm:"size:16"`
}

func (u *UserM) TableName() string {
	return "users"
}
