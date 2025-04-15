package v1

import "time"

// 用户信息
type User struct {
	// 用户ID
	UserID string `json:"userID"`
	// 用户名称
	Username string `json:"username"`
	// 用户昵称
	Nickname string `json:"nickname"`
	// 用户电子邮箱
	Email string `json:"email"`
	// 用户手机号
	Phone string `json:"phone"`
	// 用户拥有的博客数量
	PostCount int64 `json:"postCount"`
	// 用户注册时间
	CreatedAt time.Time `json:"createdAt"`
	// 用户最后更新时间
	UpdatedAt time.Time `json:"updatedAt"`
}

// 创建用户请求
type CreateUserRequest struct {
	// 用户名称
	Username string `json:"username"`
	// 用户密码
	Password string `json:"password"`
	// 用户昵称
	Nickname *string `json:"nickname"`
	// 用户电子邮箱
	Email string `json:"email"`
	// 用户手机号
	Phone string `json:"phone"`
}

// 创建用户响应
type CreateUserResponse struct {
	// 用户ID
	UserID string `json:"userID"`
}

// 更新用户请求
type UpdateUserRequest struct {
	// 可选的用户名称
	Username *string `json:"username"`
	// 可选的用户昵称
	Nickname *string `json:"nickname"`
	// 可选的用户电子邮箱
	Email *string `json:"email"`
	// 可选的用户手机号
	Phone *string `json:"phone"`
}

// 更新用户响应
type UpdateUserResponse struct {
}

// 删除用户请求
type DeleteUserRequest struct {
}

// 删除用户响应
type DeleteUserResponse struct {
}

// 获取用户请求
type GetUserRequest struct {
}

// 获取用户响应
type GetUserResponse struct {
	// 返回的用户信息
	User *User `json:"user"`
}

// 用户列表请求
type ListUserRequest struct {
	// 偏移量
	Offset int64 `json:"offset"`
	// 每页数量
	Limit int64 `json:"limit"`
}

// 用户列表响应
type ListUserResponse struct {
	// 总用户数
	TotalCount int64 `json:"totalCount"`
	// 用户列表
	Users []*User `json:"users"`
}

// LoginRequest 表示登录请求
type LoginRequest struct {
	// username 表示用户名称
	Username string `json:"username"`
	// password 表示用户密码
	Password string `json:"password"`
}

// LoginResponse 表示登录响应
type LoginResponse struct {
	// token 表示返回的身份验证令牌
	Token string `json:"token"`
	// expireAt 表示该 token 的过期时间
	ExpireAt time.Time `json:"expireAt"`
}

// RefreshTokenRequest 表示刷新令牌的请求
type RefreshTokenRequest struct {
}

// RefreshTokenResponse 表示刷新令牌的响应
type RefreshTokenResponse struct {
	// token 表示返回的身份验证令牌
	Token string `json:"token"`
	// expireAt 表示该 token 的过期时间
	ExpireAt time.Time `json:"expireAt"`
}

// ChangePasswordRequest 表示修改密码请求
type ChangePasswordRequest struct {
	// oldPassword 表示当前密码
	OldPassword string `json:"oldPassword"`
	// newPassword 表示准备修改的新密码
	NewPassword string `json:"newPassword"`
}

// ChangePasswordResponse 表示修改密码响应
type ChangePasswordResponse struct {
}
