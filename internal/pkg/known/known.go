package known

const (
	// 定义上下文中的键, 代表请求ID
	XRequestID = "x-request-id"

	// XUserID 用来定义上下文的键，代表请求用户 ID. UserID 整个用户生命周期唯一.
	XUserID = "x-user-id"

	// MaxErrGroupConcurrency 定义了 errgroup 的最大并发任务数量.
	// 用于限制 errgroup 中同时执行的 Goroutine 数量，从而防止资源耗尽，提升程序的稳定性.
	// 根据场景需求，可以调整该值大小.
	MaxErrGroupConcurrency = 1000
)
