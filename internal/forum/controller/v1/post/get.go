package post

import (
	"go-forum/internal/pkg/core"
	"go-forum/internal/pkg/log"

	"github.com/gin-gonic/gin"
)

// Get 获取 指定的博客
func (ctrl *PostController) Get(c *gin.Context) {
	log.C(c).Infow("Get post function called.")

	// 从url获取参数 postID
	postID := c.Query("postID")

	post, err := ctrl.b.Posts().Get(c, postID)

	if err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, post)
}
