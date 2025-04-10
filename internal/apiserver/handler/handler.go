package handler

import (
	"fastgo/internal/apiserver/biz"
	"fastgo/internal/apiserver/pkg/validation"
)

// 处理博客模块请求
type Handler struct {
	// HANDLER 层依赖 BIZ 层
	biz biz.IBiz
	// 封装了一系列业务校验方法
	val *validation.Validator
}

func NewHandler(biz biz.IBiz, val *validation.Validator) *Handler {
	return &Handler{
		biz: biz,
		val: val,
	}
}
