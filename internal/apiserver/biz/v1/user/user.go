package user

import (
	"context"
	"fastgo/internal/apiserver/model"
	"fastgo/internal/apiserver/pkg/conversion"
	"fastgo/internal/apiserver/store"
	"fastgo/internal/pkg/contextx"
	"fastgo/internal/pkg/known"
	where "fastgo/pkg/store"
	"log/slog"
	"sync"

	"github.com/jinzhu/copier"
	"golang.org/x/sync/errgroup"

	apiv1 "fastgo/pkg/api/apiserver/v1"
)

// UserBiz 定义处理用户请求所需的方法.
type UserBiz interface {
	// 标准资源CRUD接口
	Create(ctx context.Context, rq *apiv1.CreateUserRequest) (*apiv1.CreateUserResponse, error)
	Update(ctx context.Context, rq *apiv1.UpdateUserRequest) (*apiv1.UpdateUserResponse, error)
	Delete(ctx context.Context, rq *apiv1.DeleteUserRequest) (*apiv1.DeleteUserResponse, error)
	Get(ctx context.Context, rq *apiv1.GetUserRequest) (*apiv1.GetUserResponse, error)
	List(ctx context.Context, rq *apiv1.ListUserRequest) (*apiv1.ListUserResponse, error)

	// 扩展接口
	UserExpansion
}

// 定义用户操作的扩展方法.
type UserExpansion interface {
}

// userBiz 是 UserBiz 接口的具体实现
type userBiz struct {
	store store.IStore
}

// 静态检验 userBiz 是否实现 UserBiz 所有方法
var _ UserBiz = (*userBiz)(nil)

// 创建一个userBiz实体
func New(store store.IStore) *userBiz {
	return &userBiz{store: store}
}

// 实现 UserBiz 接口中的 Create 方法.
func (b *userBiz) Create(ctx context.Context, rq *apiv1.CreateUserRequest) (*apiv1.CreateUserResponse, error) {
	var userModel model.User
	// 将 rq 结构体赋值到 userModel
	// `copier.Copy`通过反射, 对同名/同标签的匹配字段进行复制赋值, 并忽略不匹配字段.
	_ = copier.Copy(&userModel, rq)

	// 调用 STORE 层(UserStore)的API进行数据库操作
	if err := b.store.User().Create(ctx, &userModel); err != nil {
		return nil, err
	}

	return &apiv1.CreateUserResponse{UserID: userModel.UserID}, nil
}

// 实现 UserBiz 接口的 Update 方法.
// 对 rq 的字段判空如果不为 nil 表示 request 带有这些信息
func (b *userBiz) Update(ctx context.Context, rq *apiv1.UpdateUserRequest) (*apiv1.UpdateUserResponse, error) {
	// TODO 看懂获取这个userModel的逻辑, 查询逻辑是从哪里写入的 ?
	userModel, err := b.store.User().Get(ctx, where.T(ctx))
	if err != nil {
		return nil, err
	}

	if rq.Username != nil {
		userModel.Username = *rq.Username
	}
	if rq.Email != nil {
		userModel.Email = *rq.Email
	}
	if rq.Nickname != nil {
		userModel.Nickname = *rq.Nickname
	}
	if rq.Phone != nil {
		userModel.Phone = *rq.Phone
	}

	if err := b.store.User().Update(ctx, userModel); err != nil {
		return nil, err
	}

	return &apiv1.UpdateUserResponse{}, nil
}

// 实现 UserBiz 接口中的 List 方法.
func (b *userBiz) List(ctx context.Context, rq *apiv1.ListUserRequest) (*apiv1.ListUserResponse, error) {
	// go 中 int 是32位还是64位取决操作系统
	whr := where.P(int(rq.Offset), int(rq.Limit))
	count, userList, err := b.store.User().List(ctx, whr)
	if err != nil {
		return nil, err
	}

	// 并发安全的 map
	var m sync.Map
	// TODO eg是什么?
	eg, ctx := errgroup.WithContext(ctx)

	// 设置最大协程并发数量为常量
	// 避免过高CPU或内存占用或高I/O消耗
	eg.SetLimit(known.MaxErrGroupConcurrency)

	// 使用 goroutine 提高接口性能
	for _, user := range userList {
		// 使用`eg.Go`启动的协程会按照`eg.SetLimit`的规则执行, 当达到设置的并发数时新任务会阻塞.
		eg.Go(func() error {
			select {
			case <-ctx.Done():
				return nil
			default:
				count, _, err := b.store.Post().List(ctx, where.T(ctx))
				if err != nil {
					return err
				}

				// `internal/apiserver/pkg/conversion`集成了 STORE 层返回的数据类型与 BIZ 层使用的数据类型之间的转换实现
				converted := conversion.UserodelToUserV1(user)
				converted.PostCount = count
				// 将 converted 作为值, user.ID 作为key
				m.Store(user.ID, converted)

				return nil
			}
		})
	}

	if err := eg.Wait(); err != nil {
		slog.ErrorContext(ctx, "Failed to wait all function calls returned", "err", err)
		return nil, err
	}

	// STORE 层返回的数据是降序排列的
	// 此处将 map 中的数据按照降序排入数组
	users := make([]*apiv1.User, 0, len(userList))
	for _, item := range userList {
		// TODO 待改进, 问ds反馈不应该忽略error
		user, _ := m.Load(item.ID)
		// 断言, 判定 user 是否为 *apiv1.User 类型.
		users = append(users, user.(*apiv1.User))
	}

	slog.DebugContext(ctx, "Get users from backend storage", "count", len(users))

	return &apiv1.ListUserResponse{TotalCount: count, Users: users}, nil
}

// 实现 UserBiz 接口中的 Delete 方法.
func (b *userBiz) Delete(ctx context.Context, rq *apiv1.DeleteUserRequest) (*apiv1.DeleteUserResponse, error) {
	if err := b.store.User().Delete(ctx, where.F("userID", contextx.UserID(ctx))); err != nil {
		return nil, err
	}

	return &apiv1.DeleteUserResponse{}, nil
}

// 实现 UserBiz 接口中的 Get 方法.
func (b *userBiz) Get(ctx context.Context, rq *apiv1.GetUserRequest) (*apiv1.GetUserResponse, error) {
	userModel, err := b.store.User().Get(ctx, where.F("userID", contextx.UserID(ctx)))
	if err != nil {
		return nil, err
	}

	return &apiv1.GetUserResponse{User: conversion.UserodelToUserV1(userModel)}, nil
}
