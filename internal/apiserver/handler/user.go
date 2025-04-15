package handler

import (
	"fastgo/internal/pkg/core"
	"fastgo/internal/pkg/errorsx"
	v1 "fastgo/pkg/api/apiserver/v1"
	"github.com/gin-gonic/gin"
	"log/slog"
)

// CreateUser 创建新用户.
func (h *Handler) CreateUser(c *gin.Context) {
	slog.Info("调用创建用户功能...")

	var rq v1.CreateUserRequest
	// `c.ShouldBindJSON`是gin框架提供的一个方法
	// 该方法可以将http请求体里的JSON数据解析到指定结构体中
	if err := c.ShouldBindJSON(&rq); err != nil {
		core.WriteResponse(c, errorsx.ErrBind, nil)
		return
	}

	// gin.Context 是 Gin 框架特有的上下文对象，它提供了许多处理 HTTP 请求和响应的方法，让开发者能够更方便地编写 Web 应用。
	// gin.Context.Request.Context 是 Go 标准库 net/http 中 http.Request 的 Context，主要用于管理请求的生命周期、传递请求范围内的数据以及处理超时和取消操作。
	if err := h.val.ValidateCreateUserRequest(c.Request.Context(), &rq); err != nil {
		core.WriteResponse(c, errorsx.ErrInvalidArgument.WithMessage(err.Error()), nil)
		return
	}

	// 执行具体业务逻辑
	resp, err := h.biz.UserV1().Create(c.Request.Context(), &rq)
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, resp)
}

// Login 用户登录并返回 Token.
func (h *Handler) Login(c *gin.Context) {
	slog.Info("调用用户登录功能...")

	var rq v1.LoginRequest
	if err := c.ShouldBindJSON(&rq); err != nil {
		core.WriteResponse(c, errorsx.ErrBind, nil)
		return
	}

	// 校验用户输入的登录信息
	// 输入: username, password
	if err := h.val.ValidateLoginRequest(c.Request.Context(), &rq); err != nil {
		core.WriteResponse(c, errorsx.ErrInvalidArgument, nil)
		return
	}

	resp, err := h.biz.UserV1().Login(c.Request.Context(), &rq)
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, resp)
}

// RefreshToken 刷新 JWT Token.
func (h *Handler) RefreshToken(c *gin.Context) {
	slog.Info("调用刷新token功能")

	var rq v1.RefreshTokenRequest
	if err := c.ShouldBindJSON(&rq); err != nil {
		core.WriteResponse(c, errorsx.ErrBind, nil)
		return
	}

	// 校验待补充

	resp, err := h.biz.UserV1().RefreshToken(c.Request.Context(), &rq)
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, resp)
}

func (h *Handler) ChangePassword(c *gin.Context) {
	slog.Info("调用修改密码功能")

	var rq v1.ChangePasswordRequest
	if err := c.ShouldBindJSON(&rq); err != nil {
		core.WriteResponse(c, errorsx.ErrPasswordInvalid, nil)
		return
	}

	// 校验新旧密码有效性
	if err := h.val.ValidateChangePasswordRequest(c.Request.Context(), &rq); err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	resp, err := h.biz.UserV1().ChangePassword(c.Request.Context(), &rq)
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, resp)
}
