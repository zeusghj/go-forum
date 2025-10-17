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
	CommentCounts(ctx context.Context, postIDs []uint) ([]*model.CommentCountResult, error)
	CommentCount(ctx context.Context, postID uint) (int64, error)
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
// Offset(-1).Limit(-1)操作，是 GORM 中一个用于重置或清除之前设置的分页条件的技巧，
// 目的是为了确保后续的 Count(&count)查询能够正确统计总记录数，而不是统计分页后的记录数
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
	err = p.db.Where("post_id = ?", postID).Order("created_at asc").
		Find(&ret).
		Count(&count).
		Error

	return
}

// 批量查询每条帖子对应的评论数
func (p *posts) CommentCounts(ctx context.Context, postIDs []uint) (ret []*model.CommentCountResult, err error) {
	err = p.db.Model(&model.CommentM{}).Select("post_id, count(*) as count").
		Where("post_id IN ?", postIDs).
		Group("post_id").
		Find(&ret).Error

	return
}

// 查单条帖子对应的评论数
func (p *posts) CommentCount(ctx context.Context, postID uint) (count int64, err error) {
	err = p.db.Model(&model.CommentM{}).Where("post_id = ?", postID).Count(&count).Error

	return
}
