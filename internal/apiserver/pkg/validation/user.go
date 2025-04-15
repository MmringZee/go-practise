package validation

import (
	"context"
	"errors"
	v1 "fastgo/pkg/api/apiserver/v1"
)

// ValidateCreateUserRequest 用于校验创建用户请求的输入有效性.
// 该函数包括对用户名、密码、昵称、email、手机号的校验.
func (v *Validator) ValidateCreateUserRequest(ctx context.Context, rq *v1.CreateUserRequest) error {
	// 验证用户名
	if rq.Username == "" {
		return errors.New("Username cannot be empty")
	}
	if len(rq.Username) < 4 || len(rq.Username) > 32 {
		return errors.New("Username must be between 4 and 32 characters")
	}

	// 验证密码
	if rq.Password == "" {
		return errors.New("Password cannot be empty")
	}
	if len(rq.Password) < 8 || len(rq.Password) > 64 {
		return errors.New("Password must be between 8 and 64 characters")
	}

	// 验证昵称
	if rq.Nickname != nil && *rq.Nickname != "" {
		if len(*rq.Nickname) > 32 {
			return errors.New("Nickname cannot exceed 32 characters")
		}
	}

	// 验证email
	if rq.Email == "" {
		return errors.New("Email cannot be empty")
	}

	// 验证手机号
	if rq.Phone == "" {
		return errors.New("Phone number cannot be empty")
	}

	return nil
}

// ValidateUpdateUserRequest 用于校验修改用户信息请求的输入有效性.
// 待补充...
func (v *Validator) ValidateUpdateUserRequest(ctx context.Context, rq *v1.UpdateUserRequest) error {
	return nil
}

// ValidateLoginRequest 用于校验登录请求的输入有效性.
// 对用户名和密码进行校验.
func (v *Validator) ValidateLoginRequest(ctx context.Context, rq *v1.LoginRequest) error {
	// 验证用户名
	if rq.Username == "" {
		return errors.New("Username cannot be empty")
	}
	if len(rq.Username) < 4 || len(rq.Username) > 32 {
		return errors.New("Username must be between 4 and 32 characters")
	}

	// 验证密码
	if rq.Password == "" {
		return errors.New("Password cannot be empty")
	}
	if len(rq.Password) < 8 || len(rq.Password) > 64 {
		return errors.New("Password must be between 8 and 64 characters")
	}

	return nil
}

// ValidateChangePasswordRequest 用于校验修改密码请求的密码有效性.
// 对旧密码和新密码进行校验.
func (v *Validator) ValidateChangePasswordRequest(ctx context.Context, rq *v1.ChangePasswordRequest) error {
	// 验证旧密码有效性
	if rq.OldPassword == "" {
		return errors.New("Password cannot be empty")
	}
	if len(rq.OldPassword) < 8 || len(rq.OldPassword) > 64 {
		return errors.New("Password must be between 8 and 64 characters")
	}
	// 验证新密码有效性
	if rq.NewPassword == "" {
		return errors.New("Password cannot be empty")
	}
	if len(rq.NewPassword) < 8 || len(rq.NewPassword) > 64 {
		return errors.New("Password must be between 8 and 64 characters")
	}
	// 验证新旧密码不相同
	if rq.OldPassword == rq.NewPassword {
		return errors.New("新旧密码不应该相同")
	}

	return nil
}
