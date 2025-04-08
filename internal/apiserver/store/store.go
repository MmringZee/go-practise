package store

import (
	"context"
	"fastgo/pkg/store"
	"gorm.io/gorm"
	"sync"
)

var (
	once = sync.Once{}
	// 全局变量，方便其它包直接调用已初始化好的 datastore 实例.
	// 在Go最佳实践中, 尽量减少使用包级变量, 因为其状态难以感知, 会增加维护的复杂度.
	S *datastore
)

// IStore 定义了 Store 层需要实现的方法.
type IStore interface {
	// 返回 Store 层的 *gorm.DB 实例，在少数场景下会被用到.
	DB(ctx context.Context, wheres ...where.Where) *gorm.DB
	TX(ctx context.Context, fn func(ctx context.Context) error) error

	User() UserStore
	Post() PostStore
}

// transactionKey 用于在 context.Context 中存储事务上下文的键.
// 一个空结构, 类似于Java用一个Object类作为锁的实体
type transactionKey struct{}

// datastore 是 IStore 的具体实现.
type datastore struct {
	core *gorm.DB

	// 可以根据需要添加其他数据库实例
	// fake *gorm.DB
}

// 确保 datastore 实现了 IStore 接口.
var _ IStore = (*datastore)(nil)

// NewStore 创建一个 IStore 类型的实例.
func NewStore(db *gorm.DB) *datastore {
	// 确保 S 只被初始化一次
	once.Do(func() {
		S = &datastore{db}
	})
	return S
}

// DB 根据传入的条件（wheres）对数据库实例进行筛选.
// 如果未传入任何条件，则返回上下文中的数据库实例（事务实例或核心数据库实例）.
func (store *datastore) DB(ctx context.Context, wheres ...where.Where) *gorm.DB {
	db := store.core
	// 从上下文中提取事务实例
	// ctx.Value()取值前需要先使用ctx.WithValue()设置键值
	// 提取事务实例后通过`.(*gorm.DB)`对获取的实例进行判定
	if tx, ok := ctx.Value(transactionKey{}).(*gorm.DB); ok {
		db = tx
	}

	// 遍历所有传入的条件并逐一叠加到数据库查询对象上
	for _, whr := range wheres {
		db = whr.Where(db)
	}

	return db
}

// TX 返回一个新的事务实例.
// TX方法将`*gorm.DB`类型实例注入context
// nolint: fatcontext
func (store *datastore) TX(ctx context.Context, fn func(ctx context.Context) error) error {
	// *gorm.DB.Transcation方法会自动:1.开始事务 2.根据返回值提交/回滚 3.处理panic(异常时回滚)
	return store.core.WithContext(ctx).Transaction(
		func(tx *gorm.DB) error {
			// 将GORM的事务对象*gorm.DB存入context，使业务代码可以通过context获取当前事务对象。
			// 使用空结构体transactionKey{}作为键，这是Go的惯用法:1.无内存开销,空结构体不占内存;2.保证键的唯一性
			ctx = context.WithValue(ctx, transactionKey{}, tx)
			return fn(ctx)
		},
	)
}

// User 返回一个实现了 UserStore 接口的实例.
func (store *datastore) User() UserStore {
	return newUserStore(store)
}

// Post 返回一个实现了 PostStore 接口的实例.
func (store *datastore) Post() PostStore {
	return newPostStore(store)
}
