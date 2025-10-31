package user

import (
	"go-forum/internal/pkg/core"
	"go-forum/internal/pkg/errno"
	"go-forum/internal/pkg/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 获取某用户的详细信息
func (ctrl *UserController) GetUser(c *gin.Context) {
	// 读取查询参数 ?id=xxx
	idStr := c.Query("id")

	if idStr == "" {
		// 如果没传 id，说明是获取自己的信息
		userM := c.MustGet("user").(model.UserM) // 我的理解是经过了认证中间件这里一定是有值的
		// userID, _ := c.Value(known.XUserIDKey).(uint)
		user, err := ctrl.b.Users().GetUser(c, userM.ID)
		if err != nil {
			core.WriteResponse(c, errno.ErrUserNotFound, nil)
			return
		}
		core.WriteResponse(c, nil, user)
		return
	}

	// 否则是通过 id 查询其他用户
	id, err := strconv.Atoi(idStr)
	if err != nil {
		core.WriteResponse(c, errno.ErrInvalidParameter, nil)
		return
	}

	user, err := ctrl.b.Users().GetUser(c, uint(id))
	if err != nil {
		core.WriteResponse(c, errno.ErrUserNotFound, nil)
		return
	}

	core.WriteResponse(c, nil, user)
}
