package post

import (
	"go-forum/internal/pkg/core"
	"go-forum/internal/pkg/errno"
	"go-forum/internal/pkg/known"
	"go-forum/internal/pkg/log"
	v1 "go-forum/pkg/api/forum/v1"

	"github.com/gin-gonic/gin"
)

// AddComment 对帖子添加一条评论
func (ctrl *PostController) AddComment(c *gin.Context) {
	log.C(c).Infow("Comment-Add function called")

	// 从上下文获取 user_id （JWT中间件存的）
	userIDValue, exists := c.Get(known.XUserIDKey)
	if !exists {
		core.WriteResponse(c, errno.ErrTokenInvalid, nil)
		return
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		core.WriteResponse(c, errno.ErrTokenInvalid.SetMessage("用户ID解析失败"), nil)
		return
	}

	var r v1.CreateCommentRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, errno.ErrInvalidParameter.SetMessage(err.Error()), nil)

		return
	}

	if err := ctrl.b.Posts().AddComment(c, userID, &r); err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, gin.H{"code": 200, "data": "评论成功"})
}
