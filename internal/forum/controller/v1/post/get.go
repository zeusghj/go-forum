package post

import (
	"go-forum/internal/pkg/core"
	"go-forum/internal/pkg/errno"
	"go-forum/internal/pkg/log"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Get 获取 指定的博客
func (ctrl *PostController) Get(c *gin.Context) {
	log.C(c).Infow("Get post function called.")

	// 从url获取参数 postID
	postIDStr := c.Query("id")

	if postIDStr == "" {
		core.WriteResponse(c, errno.ErrInvalidParameter.SetMessage("未填写post id"), nil)
		return
	}

	// 使用 base=10, bitSize=64 来解析为uint64，然后转换为uint 将字符串参数解析为无符号整数
	postID, err := strconv.ParseUint(postIDStr, 10, 64)
	if err != nil {
		// 这里可以更精细地处理错误类型
		core.WriteResponse(c, errno.ErrInvalidParameter.SetMessage("post id 必须是一个非负整数"), nil)
		return
	}

	post, err := ctrl.b.Posts().Get(c, uint(postID))

	if err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, post)
}
