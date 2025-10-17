package store

import (
	"context"
	"go-forum/internal/pkg/model"

	"gorm.io/gorm"
)

// UserStore 定义了 user 模块在 store 层所实现的方法.
type UserStore interface {
	Create(ctx context.Context, user *model.UserM) error
	GetByID(ctx context.Context, userID uint) (*model.UserM, error)
	GetByUsername(ctx context.Context, username string) (*model.UserM, error)
	Update(ctx context.Context, user *model.UserM) error
	GetUsers(ctx context.Context, userIDs []uint) ([]*model.UserM, error)
}

// UserStore 接口的实现.
type users struct {
	db *gorm.DB
}

// 确保 users 实现了 UserStore 接口.
var _ UserStore = (*users)(nil)

func newUsers(db *gorm.DB) *users {
	return &users{db}
}

// Create 插入一条 user 记录.
func (u *users) Create(ctx context.Context, user *model.UserM) error {
	return u.db.Create(&user).Error
}

// Get 根据用户id查询指定 user 的数据库记录.
func (u *users) GetByID(ctx context.Context, userID uint) (*model.UserM, error) {
	var user model.UserM
	if err := u.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// Get 根据用户名查询指定 user 的数据库记录.
func (u *users) GetByUsername(ctx context.Context, username string) (*model.UserM, error) {
	var user model.UserM
	if err := u.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// 批量查用户信息 参数 userIDs
func (u *users) GetUsers(ctx context.Context, userIDs []uint) (ret []*model.UserM, err error) {

	err = u.db.Where("id IN ?", userIDs).Find(&ret).Error

	return
}

// Update 更新一条 user 数据库记录.
func (u *users) Update(ctx context.Context, user *model.UserM) error {
	return u.db.Save(user).Error
}
