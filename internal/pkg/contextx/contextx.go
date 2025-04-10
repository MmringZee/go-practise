package contextx

import "context"

// 定义用于上下文的键
type (
	// 仅将这个结构体作为一个key
	requestIDKey struct{}
	// userIDKey 定义用户 ID 的上下文键.
	userIDKey struct{}
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

// 将用户 ID 存放到上下文中.
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey{}, userID)
}

// 从上下文中提取用户 ID.
func UserID(ctx context.Context) string {
	userID, _ := ctx.Value(userIDKey{}).(string)
	return userID
}
