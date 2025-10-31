package model

import (
	"go-forum/pkg/auth"

	"gorm.io/gorm"
)

type UserM struct {
	gorm.Model
	Username string `gorm:"unique;size:255;not null"`
	Password string `gorm:"size:255;not null"`
	Nickname string `gorm:"size:30"`
	Email    string `gorm:"size:256"`
	Phone    string `gorm:"size:16"`
	Role     string `gorm:"size:50;default:'user'"`
}

func (u *UserM) TableName() string {
	return "users"
}

// BeforeCreate 在创建数据库记录之前加密明文密码.
func (u *UserM) BeforeCreate(tx *gorm.DB) (err error) {
	// Encrypt the user password.
	u.Password, err = auth.Encrypt(u.Password)
	if err != nil {
		return err
	}

	return nil
}
