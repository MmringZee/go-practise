package validation

import (
	"context"
	"errors"
	v1 "fastgo/pkg/api/apiserver/v1"
)

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

func (v *Validator) ValidateUpdateUserRequest(ctx context.Context, rq *v1.UpdateUserRequest) error {
	return nil
}
