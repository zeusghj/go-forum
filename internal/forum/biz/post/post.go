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
	AddComment(ctx context.Context, userID uint, r *v1.CreateCommentRequest) error
	CommentList(ctx context.Context, postID uint) (*v1.ListCommentResponse, error)
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

	// 查询作者名称 (可选)
	authorName := ""
	user, err := p.ds.Users().GetByID(ctx, postM.UserID)
	if err == nil {
		authorName = user.Username
	}

	// 查询评论数
	commentCount, err := p.ds.Posts().CommentCount(ctx, postM.ID)
	if err != nil {
		log.Errorw("查询帖子的评论数失败", "postID", postM.ID)
		commentCount = 0
	}

	var postR = v1.GetPostResponse{
		ID:           postM.ID,
		Title:        postM.Title,
		Content:      postM.Content,
		UserID:       postM.UserID,
		Username:     authorName, // 若为空，前端可处理显示“匿名”或 id
		CreatedAt:    postM.CreatedAt,
		CommentCount: commentCount,
	}

	return &postR, nil
}

// List 是 PostBiz 接口中 `List` 方法的实现.
func (b *postBiz) List(ctx context.Context, offset, limit int) (*v1.ListPostResponse, error) {
	count, list, err := b.ds.Posts().List(ctx, offset, limit)
	if err != nil {
		log.C(ctx).Errorw("Failed to list posts from storage", "err", err)
		return nil, err
	}

	// 收集 post ids 和 user ids
	postIDs := make([]uint, 0, len(list))
	userIDSet := make(map[uint]struct{})
	userIDs := make([]uint, 0, len(list))
	for _, p := range list {
		postIDs = append(postIDs, p.ID)
		if _, ok := userIDSet[p.UserID]; !ok {
			userIDSet[p.UserID] = struct{}{}
			userIDs = append(userIDs, p.UserID)
		}
	}

	// 查询用户信息 （批量）
	userMap := map[uint]string{}
	if len(userIDs) > 0 {
		users, err := b.ds.Users().GetUsers(ctx, userIDs)
		if err != nil {
			log.Errorw("批量查询用户信息出错")
			return nil, err
		}
		for _, u := range users {
			userMap[u.ID] = u.Username
		}
	}

	// 查询每个帖子的评论数量
	commentCountMap := map[uint]int64{}
	if len(postIDs) > 0 {
		list, err := b.ds.Posts().CommentCounts(ctx, postIDs)
		if err != nil {
			log.Errorw("查询每个帖子的评论数量出错")
			return nil, err
		}
		for _, ct := range list {
			commentCountMap[ct.PostID] = ct.Count
		}
	}

	// 构造响应

	posts := make([]*v1.PostInfo, 0, len(list))
	for _, p := range list {
		posts = append(posts, &v1.PostInfo{
			ID:           p.ID,
			Title:        p.Title,
			Content:      p.Content,
			UserID:       p.UserID,
			Username:     userMap[p.UserID], // 若为空，前端可处理显示“匿名”或 id
			CreatedAt:    p.CreatedAt,
			CommentCount: commentCountMap[p.ID],
		})
	}

	return &v1.ListPostResponse{TotalCount: count, Posts: posts}, nil
}

// AddComment 是 PostBiz 接口中 `AddComment` 方法的实现.
func (b *postBiz) AddComment(ctx context.Context, userID uint, r *v1.CreateCommentRequest) error {
	var commentM model.CommentM
	_ = copier.Copy(&commentM, r)

	commentM.UserID = userID

	if err := b.ds.Posts().AddComment(ctx, &commentM); err != nil {
		return err
	}

	return nil
}

func (b *postBiz) CommentList(ctx context.Context, postID uint) (*v1.ListCommentResponse, error) {
	count, list, err := b.ds.Posts().CommentList(ctx, postID)
	if err != nil {
		log.C(ctx).Errorw("Failed to list comments from storage", "err", err)
		return nil, err
	}

	// 收集评论用户的 user ids
	userIDSet := map[uint]struct{}{}
	userIDs := make([]uint, 0, len(list))
	for _, cm := range list {
		if _, ok := userIDSet[cm.UserID]; !ok {
			userIDSet[cm.UserID] = struct{}{}
			userIDs = append(userIDs, cm.UserID)
		}
	}

	// 批量查询评论用户的用户名
	userMap := map[uint]string{}
	if len(userIDs) > 0 {
		users, err := b.ds.Users().GetUsers(ctx, userIDs)
		if err != nil {
			log.Errorw("批量查询用户信息出错")
			return nil, err
		}
		for _, u := range users {
			userMap[u.ID] = u.Username
		}
	}

	comments := make([]*v1.CommentInfo, 0, len(list))
	for _, cm := range list {
		comments = append(comments, &v1.CommentInfo{
			ID:        cm.ID,
			PostID:    cm.PostID,
			UserID:    cm.UserID,
			Username:  userMap[cm.UserID],
			Content:   cm.Content,
			CreatedAt: cm.CreatedAt,
			UpdatedAt: cm.UpdatedAt,
		})
	}

	return &v1.ListCommentResponse{TotalCount: count, Posts: comments}, nil
}
