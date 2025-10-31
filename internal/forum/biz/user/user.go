package user

import (
	"context"
	"errors"
	"regexp"

	"go-forum/internal/forum/store"
	"go-forum/internal/pkg/errno"
	"go-forum/internal/pkg/model"
	v1 "go-forum/pkg/api/forum/v1"
	"go-forum/pkg/auth"
	"go-forum/pkg/token"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

// UserBiz 定义了 user 模块在 biz 层所实现的方法.
type UserBiz interface {
	Create(ctx context.Context, r *v1.CreateUserRequest) error
	Login(ctx context.Context, r *v1.LoginRequest) (*v1.LoginResponse, error)
	GetUser(ctx context.Context, userID uint) (*v1.GetUserResponse, error)
	ChangePassword(ctx context.Context, userId uint, r *v1.ChangePasswordRequest) error
}

// UserBiz 接口的实现.
type userBiz struct {
	ds store.IStore
}

// 确保 userBiz 实现了 UserBiz 接口.
var _ UserBiz = (*userBiz)(nil)

// New 创建一个实现了 UserBiz 接口的实例.
func New(ds store.IStore) *userBiz {
	return &userBiz{ds}
}

// Create 是 UserBiz 接口中 `Create` 方法的实现.
func (b *userBiz) Create(ctx context.Context, r *v1.CreateUserRequest) error {
	var userM model.UserM
	_ = copier.Copy(&userM, r)

	if err := b.ds.Users().Create(ctx, &userM); err != nil {
		if match, _ := regexp.MatchString("Duplicate entry '.*' for key '.+'", err.Error()); match {
			return errno.ErrUserAlreadyExist
		}

		return err
	}

	return nil
}

// Login 是 UserBiz 接口中 `Login` 方法的实现.
func (u *userBiz) Login(ctx context.Context, r *v1.LoginRequest) (*v1.LoginResponse, error) {
	// 获取登录用户的所有信息
	user, err := u.ds.Users().GetByUsername(ctx, r.Username)
	if err != nil {
		return nil, errno.ErrUserNotFound
	}

	// 对比传入的明文密码和数据库中已加密过的密码是否匹配
	if err := auth.Compare(user.Password, r.Password); err != nil {
		return nil, errno.ErrPasswordIncorrect
	}

	// 如果匹配成功，说明登录成功，签发 token 并返回
	t, err := token.Sign(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, errno.ErrSignToken
	}

	// 如果匹配成功，说明登录成功， 签发 token 并返回
	return &v1.LoginResponse{Token: t}, nil
}

// ChangePassword 是 UserBiz 接口中 `ChangePassword` 方法的实现.
func (u *userBiz) ChangePassword(ctx context.Context, userID uint, r *v1.ChangePasswordRequest) error {
	user, err := u.ds.Users().GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if err := auth.Compare(user.Password, r.OldPassword); err != nil {
		return errno.ErrPasswordIncorrect
	}

	user.Password, _ = auth.Encrypt(r.NewPassword)
	if err := u.ds.Users().Update(ctx, user); err != nil {
		return err
	}

	return nil
}

// 获取用户详情
func (u *userBiz) GetUser(ctx context.Context, userID uint) (*v1.GetUserResponse, error) {
	user, err := u.ds.Users().GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.ErrUserNotFound
		}
		return nil, err
	}

	var resp v1.GetUserResponse
	_ = copier.CopyWithOption(&resp, user, copier.Option{IgnoreEmpty: true, DeepCopy: true})

	return &resp, nil
}
