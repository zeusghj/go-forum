package post

import (
	"go-forum/internal/forum/biz"
	"go-forum/internal/forum/store"
)

// PostController 是 post 模块在 Controller 层的实现，用来处理 post 模块的请求.
type PostController struct {
	b biz.IBiz
}

// New 创建一个 post controller
func New(ds store.IStore) *PostController {
	return &PostController{b: biz.NewBiz(ds)}
}
