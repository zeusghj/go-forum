package user

import (
	"go-forum/internal/pkg/core"
	"go-forum/internal/pkg/errno"
	"go-forum/internal/pkg/known"
	"go-forum/internal/pkg/log"
	v1 "go-forum/pkg/api/forum/v1"

	"github.com/gin-gonic/gin"
)

// ChangePassword 用来修改指定用户的密码.
func (ctrl *UserController) ChangePassword(c *gin.Context) {
	log.C(c).Infow("ChangePassword function called")

	userId := c.Value(known.XUsernameKey)
	username, ok := userId.(string)

	if !ok || username == "" {
		// 当前登录的用户有问题
		core.WriteResponse(c, errno.ErrTokenInvalid, nil)

		return
	}

	var r v1.ChangePasswordRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, errno.ErrInvalidParameter.SetMessage(err.Error()), nil)

		return
	}

	if err := ctrl.b.Users().ChangePassword(c, username, &r); err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, gin.H{"message": "修改密码成功"})
}
