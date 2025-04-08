// 整个 UserStore 接口中的方法实现较为简洁，仅直接对数据库记录进行增删改查操作，并未封装任何业务逻辑。
// 在 Go 项目开发中，不少开发者会在 Store 层封装业务代码，为了实现不同的查询条件，会在 Store 层封装很多查询类方法.
// 例如：ListUser、ListUserByName、ListUserByID 等。这都会使 Store 层代码变得臃肿，难以维护。
// 其实 Store 层只需要对数据库记录进行简单的增删改查即可。对插入数据或查询数据的处理可以放在业务逻辑层。对查询条件的定制，可以通过提供灵活的查询参数来实现。
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

// UserStore 定义了 user 模块在 store 层实现的方法.
type UserStore interface {
	Create(ctx context.Context, obj *model.User) error
	Update(ctx context.Context, obj *model.User) error
	Delete(ctx context.Context, opts *where.Options) error
	Get(ctx context.Context, opts *where.Options) (*model.User, error)
	List(ctx context.Context, opts *where.Options) (int64, []*model.User, error)

	UserExpansion
}

// UserExpansion 定义了用户操作的附加方法.
type UserExpansion interface {
}

// userStore 是 UserStore 接口的实现.
// Go 一般通过组合的方式完成"继承" / 通过实现方法的方式完成"实现接口"
type userStore struct {
	// 组合了 datastore 结构体
	store *datastore
}

// Go 语言中用于静态接口实现检查的惯用模式, 通过编译时的类型断言来确保结构体完整实现了接口的所有方法
// 创建一个 nil 的 userStore 指针
var _ UserStore = (*userStore)(nil)

// newUserStore 创建 userStore 的实例.
func newUserStore(store *datastore) *userStore {
	return &userStore{store: store}
}

// Create 插入一条用户记录.
func (s *userStore) Create(ctx context.Context, obj *model.User) error {
	// 调用`s.store.DB(ctx)`尝试从context中获取事务, 若没有事务则获取`*gorm.DB`类型的实例
	// 调用`*gorm.DB`提供的`Create`方法进行数据库插入记录
	if err := s.store.DB(ctx).Create(&obj).Error; err != nil {
		slog.Error("Failed to insert user into database", "err", err, "user", obj)
		// 项目`internal/pkg/errorsx`对DB错误进行封装, 防止直接输出未过滤的敏感信息
		return errorsx.ErrDBWrite.WithMessage(err.Error())
	}
	return nil
}

// Delete 根据条件删除用户记录.
func (s *userStore) Delete(ctx context.Context, opts *where.Options) error {
	err := s.store.DB(ctx).Delete(new(model.User)).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		slog.Error("Failed to delete user from database", "err", err, "conditions", opts)
		return errorsx.ErrDBWrite.WithMessage(err.Error())
	}
	return nil
}

// List 返回用户列表和总数.
// nolint: nonamedreturns
func (s *userStore) List(ctx context.Context, opts *where.Options) (count int64, ret []*model.User, err error) {
	// 通过`s.store.DB`的可变入参传入查询条件
	// 后续表示 : 按数据库字段`id`降序排列、
	err = s.store.DB(ctx, opts).Order("id desc").Find(&ret).Offset(-1).Limit(-1).Count(&count).Error
	if err != nil {
		slog.Error("Failed to list users from database", "err", err, "conditions", opts)
		err = errorsx.ErrDBRead.WithMessage(err.Error())
	}
	return
}

// Update 更新用户数据库记录.
func (s *userStore) Update(ctx context.Context, obj *model.User) error {
	if err := s.store.DB(ctx).Save(obj).Error; err != nil {
		slog.Error("Failed to update user in database", "err", err, "user", obj)
		return errorsx.ErrDBWrite.WithMessage(err.Error())
	}
	return nil
}

// Get 根据条件查询用户记录.
func (s *userStore) Get(ctx context.Context, opts *where.Options) (*model.User, error) {
	var obj model.User
	if err := s.store.DB(ctx, opts).First(&obj).Error; err != nil {
		slog.Error("Failed to retrieve user from database", "err", err, "conditions", opts)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorsx.ErrUserNotFound
		}
		return nil, errorsx.ErrDBRead.WithMessage(err.Error())
	}
	return &obj, nil
}
