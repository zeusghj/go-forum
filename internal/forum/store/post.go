package store

import (
	"context"
	"go-forum/internal/pkg/model"

	"gorm.io/gorm"
)

// PostStore 定义了 post 模块在 store 层所实现的方法.
type PostStore interface {
	Create(ctx context.Context, post *model.PostM) error
	Get(ctx context.Context, postID uint) (*model.PostM, error)
	List(ctx context.Context, offset, limit int) (int64, []*model.PostM, error)
	AddComment(ctx context.Context, comment *model.CommentM) error
	CommentList(ctx context.Context, postID uint) (int64, []*model.CommentM, error)
}

// PostStore 接口的实现
type posts struct {
	db *gorm.DB
}

// 确保 posts 实现了 PostStore 接口
var _ PostStore = (*posts)(nil)

func newPosts(db *gorm.DB) *posts {
	return &posts{db}
}

// Create 插入一条 post 记录
func (p *posts) Create(ctx context.Context, post *model.PostM) error {
	return p.db.Create(&post).Error
}

// Get 根据 postID 查询指定 post 的数据库记录.
func (p *posts) Get(ctx context.Context, postID uint) (*model.PostM, error) {
	var post model.PostM
	if err := p.db.Where("id = ?", postID).First(&post).Error; err != nil {
		return nil, err
	}

	return &post, nil
}

// List 根据 offset 和 limit 返回 post 列表.
func (p *posts) List(ctx context.Context, offset, limit int) (count int64, ret []*model.PostM, err error) {
	err = p.db.Offset(offset).Limit(defaultLimit(limit)).Order("id desc").
		Find(&ret).
		Offset(-1).
		Limit(-1).
		Count(&count).
		Error

	return
}

// AddComment 对postID 添加一条评论
func (p *posts) AddComment(ctx context.Context, comment *model.CommentM) error {
	return p.db.Create(&comment).Error
}

// CommentList  返回 comment 列表.
func (p *posts) CommentList(ctx context.Context, postID uint) (count int64, ret []*model.CommentM, err error) {
	err = p.db.Where("post_id = ?", postID).Order("id desc").
		Find(&ret).
		Count(&count).
		Error

	return
}
