package store

import (
	"context"
	"errors"
	"fastgo/internal/apiserver/model"
	"fastgo/internal/pkg/errorsx"
	where "fastgo/pkg/store"
	"gorm.io/gorm"
	"log/slog"
)

// PostStore 定义了 post 模块在 store 层实现的方法.
type PostStore interface {
	Create(ctx context.Context, obj *model.Post) error
	Update(ctx context.Context, obj *model.Post) error
	Delete(ctx context.Context, opts *where.Options) error
	Get(ctx context.Context, opts *where.Options) (*model.Post, error)
	List(ctx context.Context, opts *where.Options) (int64, []*model.Post, error)

	PostExpansion
}

// PostExpansion 定义了用户操作的附加方法.
type PostExpansion interface {
}

type postStore struct {
	store *datastore
}

var _ PostStore = (*postStore)(nil)

// newPostStore 创建 postStore 的实例.
func newPostStore(store *datastore) *postStore {
	return &postStore{store: store}
}

// Create 插入一条用户记录.
func (s *postStore) Create(ctx context.Context, obj *model.Post) error {
	// 调用`s.store.DB(ctx)`尝试从context中获取事务, 若没有事务则获取`*gorm.DB`类型的实例
	// 调用`*gorm.DB`提供的`Create`方法进行数据库插入记录
	if err := s.store.DB(ctx).Create(&obj).Error; err != nil {
		slog.Error("Failed to insert post into database", "err", err, "post", obj)
		// 项目`internal/pkg/errorsx`对DB错误进行封装, 防止直接输出未过滤的敏感信息
		return errorsx.ErrDBWrite.WithMessage(err.Error())
	}
	return nil
}

// Delete 根据条件删除用户记录.
func (s *postStore) Delete(ctx context.Context, opts *where.Options) error {
	err := s.store.DB(ctx).Delete(new(model.Post)).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		slog.Error("Failed to delete post from database", "err", err, "conditions", opts)
		return errorsx.ErrDBWrite.WithMessage(err.Error())
	}
	return nil
}

// List 返回用户列表和总数.
// nolint: nonamedreturns
func (s *postStore) List(ctx context.Context, opts *where.Options) (count int64, ret []*model.Post, err error) {
	// 通过`s.store.DB`的可变入参传入查询条件
	// 后续表示 : 按数据库字段`id`降序排列、
	err = s.store.DB(ctx, opts).Order("id desc").Find(&ret).Offset(-1).Limit(-1).Count(&count).Error
	if err != nil {
		slog.Error("Failed to list posts from database", "err", err, "conditions", opts)
		err = errorsx.ErrDBRead.WithMessage(err.Error())
	}
	return
}

// Update 更新帖子数据库记录.
func (s *postStore) Update(ctx context.Context, obj *model.Post) error {
	if err := s.store.DB(ctx).Save(obj).Error; err != nil {
		slog.Error("Failed to update post in database", "err", err, "post", obj)
		return errorsx.ErrDBWrite.WithMessage(err.Error())
	}
	return nil
}

// Get 根据条件查询帖子记录.
func (s *postStore) Get(ctx context.Context, opts *where.Options) (*model.Post, error) {
	var obj model.Post
	if err := s.store.DB(ctx, opts).First(&obj).Error; err != nil {
		slog.Error("Failed to retrieve post from database", "err", err, "conditions", opts)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorsx.ErrPostNotFound
		}
		return nil, errorsx.ErrDBRead.WithMessage(err.Error())
	}
	return &obj, nil
}
