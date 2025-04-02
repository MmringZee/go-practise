// 这段代码展示了 Go 语言中常用的选项模式(Options Pattern)实现
// 主要用于灵活配置两种不同类型的选项：CodeOptions 和 SonyflakeOptions。
package id

import "time"

type CodeOptions struct {
	// 可用字符集合
	chars []rune
	// 与字符集长度互质的数
	n1 int
	// 与代码长度互质的数
	n2 int
	// 生成的代码长度
	l int
	// 随机盐值
	salt uint64
}

// 相同类型入参和返回, 便于链式调用和嵌套使用
func getCodeOptionsOrSetDefault(options *CodeOptions) *CodeOptions {
	if options == nil {
		return &CodeOptions{
			// remove 0,1,I,O,U,Z
			// 去除易混淆字符
			chars: []rune{
				'2', '3', '4', '5', '6',
				'7', '8', '9', 'A', 'B',
				'C', 'D', 'E', 'F', 'G',
				'H', 'J', 'K', 'L', 'M',
				'N', 'P', 'Q', 'R', 'S',
				'T', 'V', 'W', 'X', 'Y',
			},
			// n1选择与字符集长度30互质的数, 减少冲突概率
			n1: 17,
			// 与代码长度l互质的数
			n2: 5,
			// 代码长度
			l: 8,
			// 随机盐值
			salt: 123567369,
		}
	}
	return options
}

// 选项配置函数, 返回一个闭包函数用于修改选项
// WithXXXXXX
// 以WithCodeChars方法为例, 该闭包方法实现了 : 传入一个arr字符数组,
func WithCodeChars(arr []rune) func(*CodeOptions) {
	return func(options *CodeOptions) {
		if len(arr) > 0 {
			getCodeOptionsOrSetDefault(options).chars = arr
		}
	}
}

func WithCodeN1(n int) func(*CodeOptions) {
	return func(options *CodeOptions) {
		getCodeOptionsOrSetDefault(options).n1 = n
	}
}

func WithCodeN2(n int) func(*CodeOptions) {
	return func(options *CodeOptions) {
		getCodeOptionsOrSetDefault(options).n2 = n
	}
}

func WithCodeL(l int) func(*CodeOptions) {
	return func(options *CodeOptions) {
		if l > 0 {
			getCodeOptionsOrSetDefault(options).l = l
		}
	}
}

func WithCodeSalt(salt uint64) func(*CodeOptions) {
	return func(options *CodeOptions) {
		if salt > 0 {
			getCodeOptionsOrSetDefault(options).salt = salt
		}
	}
}

type SonyflakeOptions struct {
	machineId uint16
	startTime time.Time
}

func getSonyflakeOptionsOrSetDefault(options *SonyflakeOptions) *SonyflakeOptions {
	if options == nil {
		return &SonyflakeOptions{
			machineId: 1,
			startTime: time.Date(2022, 10, 10, 0, 0, 0, 0, time.UTC),
		}
	}
	return options
}

func WithSonyflakeMachineId(id uint16) func(*SonyflakeOptions) {
	return func(options *SonyflakeOptions) {
		if id > 0 {
			getSonyflakeOptionsOrSetDefault(options).machineId = id
		}
	}
}

func WithSonyflakeStartTime(startTime time.Time) func(*SonyflakeOptions) {
	return func(options *SonyflakeOptions) {
		if !startTime.IsZero() {
			getSonyflakeOptionsOrSetDefault(options).startTime = startTime
		}
	}
}
