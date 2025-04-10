package handler

import (
	"fastgo/internal/pkg/core"
	"fastgo/internal/pkg/errorsx"
	v1 "fastgo/pkg/api/apiserver/v1"
	"github.com/gin-gonic/gin"
	"log/slog"
)

func (h *Handler) CreateUser(c *gin.Context) {
	slog.Info("调用创建用户功能...")

	var rq v1.CreateUserRequest
	// `c.ShouldBindJSON`是gin框架提供的一个方法
	// 该方法可以将http请求体里的JSON数据解析到指定结构体中
	if err := c.ShouldBindJSON(&rq); err != nil {
		core.WriteResponse(c, nil, errorsx.ErrBind)
		return
	}

	// gin.Context 是 Gin 框架特有的上下文对象，它提供了许多处理 HTTP 请求和响应的方法，让开发者能够更方便地编写 Web 应用。
	// gin.Context.Request.Context 是 Go 标准库 net/http 中 http.Request 的 Context，主要用于管理请求的生命周期、传递请求范围内的数据以及处理超时和取消操作。
	if err := h.val.ValidateCreateUserRequest(c.Request.Context(), &rq); err != nil {
		core.WriteResponse(c, nil, errorsx.ErrInvalidArgument.WithMessage(err.Error()))
		return
	}

	// 执行具体业务逻辑
	resp, err := h.biz.UserV1().Create(c.Request.Context(), &rq)
	if err != nil {
		core.WriteResponse(c, nil, err)
		return
	}

	core.WriteResponse(c, nil, resp)
}
