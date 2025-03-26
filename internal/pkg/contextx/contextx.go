package contextx

import "context"

// 定义用于上下文的键
type (
	// 仅将这个结构体作为一个key
	requestIDKey struct{}
)

// 将请求ID存放到上下文中
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey{}, requestID)
}

// 从上下文中提取请求ID
func RequestID(ctx context.Context) string {
	requestID, _ := ctx.Value(requestIDKey{}).(string)
	return requestID
}
