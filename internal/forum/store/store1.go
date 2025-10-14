package store

import (
	"go-forum/internal/pkg/model"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	// 修改成你自己的 MySQL 配置 用户名是 forum_user，密码是 123456 ,数据库是 go_forum
	dsn := "forum_user:123456@tcp(127.0.0.1:3306)/go_forum?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("数据库连接失败：", err)
	}

	// 自动建表
	DB.AutoMigrate(&model.UserM{}, &model.PostM{}, &model.CommentM{})
}
