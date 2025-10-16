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

	userID := c.Value(known.XUserIDKey).(uint)

	var r v1.CreateCommentRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, errno.ErrInvalidParameter.SetMessage(err.Error()), nil)

		return
	}

	if err := ctrl.b.Posts().AddComment(c, userID, &r); err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, gin.H{"code": 200, "data": "成功"})
}
