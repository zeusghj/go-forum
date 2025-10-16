package post

import (
	"context"
	"errors"
	"go-forum/internal/forum/store"
	"go-forum/internal/pkg/errno"
	"go-forum/internal/pkg/log"
	"go-forum/internal/pkg/model"
	v1 "go-forum/pkg/api/forum/v1"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

// PostBiz 定义了 post 模块在 biz 层所实现的方法
type PostBiz interface {
	Create(ctx context.Context, userID uint, r *v1.CreatePostRequest) error
	Get(ctx context.Context, postID uint) (*v1.GetPostResponse, error)
	List(ctx context.Context, offset, limit int) (*v1.ListPostResponse, error)
}

// PostBiz 接口的实现
type postBiz struct {
	ds store.IStore
}

// 确保 userBiz 实现了 UserBiz 接口.
var _ PostBiz = (*postBiz)(nil)

// New 创建一个实现了 PostBiz 接口的实例.
func New(ds store.IStore) *postBiz {
	return &postBiz{ds}
}

// Create 是 PostBiz 接口中 `Create` 方法的实现.
func (b *postBiz) Create(ctx context.Context, userID uint, r *v1.CreatePostRequest) error {
	var postM model.PostM
	_ = copier.Copy(&postM, r)

	postM.UserID = userID

	if err := b.ds.Posts().Create(ctx, &postM); err != nil {
		return err
	}

	return nil
}

// Get implements PostBiz.
func (p *postBiz) Get(ctx context.Context, postID uint) (*v1.GetPostResponse, error) {
	postM, err := p.ds.Posts().Get(ctx, postID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.ErrPostNotFound
		}

		return nil, err
	}

	var postR v1.GetPostResponse
	_ = copier.Copy(&postR, postM)

	return &postR, nil
}

// List 是 PostBiz 接口中 `List` 方法的实现.
func (b *postBiz) List(ctx context.Context, offset, limit int) (*v1.ListPostResponse, error) {
	count, list, err := b.ds.Posts().List(ctx, offset, limit)
	if err != nil {
		log.C(ctx).Errorw("Failed to list posts from storage", "err", err)
		return nil, err
	}

	posts := make([]*v1.PostInfo, 0, len(list))
	for _, item := range list {
		post := item
		posts = append(posts, &v1.PostInfo{
			Username:  "还未实现",
			ID:        post.ID,
			Title:     post.Title,
			Content:   post.Content,
			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
		})
	}

	return &v1.ListPostResponse{TotalCount: count, Posts: posts}, nil
}
