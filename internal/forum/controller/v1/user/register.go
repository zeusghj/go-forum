package user

import (
	"go-forum/internal/pkg/core"
	"go-forum/internal/pkg/errno"
	"go-forum/internal/pkg/log"
	v1 "go-forum/pkg/api/forum/v1"

	"github.com/gin-gonic/gin"
)

// const defaultMethods = "(GET)|(POST)|(PUT)|(DELETE)"

// Register 创建一个新的用户
func (ctrl *UserController) Register(c *gin.Context) {
	log.C(c).Infow("register function called")

	var r v1.CreateUserRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, errno.ErrInvalidParameter.SetMessage(err.Error()), nil)

		return
	}

	if err := ctrl.b.Users().Create(c, &r); err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	// if _, err := ctrl.a.AddNamedPolicy("p", r.Username, "/v1/user/profile", defaultMethods); err != nil {
	// 	core.WriteResponse(c, err, nil)

	// 	return
	// }

	role := "user"
	if r.Username == "root" {
		role = "admin"
	}
	// AddNamedPolicy 是用来添加 p策略（policy） 的方法，即 p, sub, obj, act 类型的规则。
	// if _, err := ctrl.a.AddNamedPolicy("g", r.Username, role); err != nil {
	// 	core.WriteResponse(c, err, nil)

	// 	return
	// }

	// AddNamedGroupingPolicy 某个用户属于某个角色
	if _, err := ctrl.a.AddNamedGroupingPolicy("g", r.Username, role); err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, gin.H{"message": "注册成功"})
}
