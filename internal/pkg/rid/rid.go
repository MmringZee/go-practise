package rid

import (
	"fastgo/internal/pkg/id"
)

const defaultABC = "abcdefghijklmnopqrstuvwxyz1234567890"

type ResourceID string

const (
	// 资源标识符
	UserID ResourceID = "user"
	PostID ResourceID = "post"
)

// 将资源标识符转换为字符串
func (rid ResourceID) String() string {
	return rid.String()
}

// 创建带前缀的唯一标识符
func (rid ResourceID) New(connter uint64) string {
	// 使用自定义选项生成唯一标识符
	// NewCode()方法传入了一系列闭包函数, 这些闭包函数指定了入参, 可以根据这些入参对NewCode
	// 需要注意, WithCodeL需要传入一个与5(CodeOptions.n2默认值)互质的数
	uniqueStr := id.NewCode(
		connter,
		id.WithCodeChars([]rune(defaultABC)),
		id.WithCodeL(6),
		id.WithCodeSalt(Salt()),
	)
	return rid.String() + "-" + uniqueStr
}
