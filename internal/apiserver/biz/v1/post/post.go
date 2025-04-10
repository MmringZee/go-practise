package post

import (
	"context"
	"fastgo/internal/apiserver/model"
	"fastgo/internal/apiserver/pkg/conversion"
	"fastgo/internal/apiserver/store"
	"fastgo/internal/pkg/contextx"
	where "fastgo/pkg/store"
	"github.com/jinzhu/copier"

	apiv1 "fastgo/pkg/api/apiserver/v1"
)

// 定义处理帖子请求所需的方法
type PostBiz interface {
	Create(ctx context.Context, rq *apiv1.CreatePostRequest) (*apiv1.CreatePostResponse, error)
	Update(ctx context.Context, rq *apiv1.UpdatePostRequest) (*apiv1.UpdatePostResponse, error)
	Delete(ctx context.Context, rq *apiv1.DeletePostRequest) (*apiv1.DeletePostResponse, error)
	Get(ctx context.Context, rq *apiv1.GetPostRequest) (*apiv1.GetPostResponse, error)
	List(ctx context.Context, rq *apiv1.ListPostRequest) (*apiv1.ListPostResponse, error)

	PostExpansion
}

// 定义额外的帖子操作方法.
type PostExpansion interface{}

// PostBiz 接口的实现.
// BIZ 层依赖 STORE 层, 通过组合方式实现依赖
type postBiz struct {
	store store.IStore
}

// 创建 postBiz 实例
func New(store store.IStore) *postBiz {
	return &postBiz{
		store: store,
	}
}

// 静态检验接口函数都已实现
var _ PostBiz = (*postBiz)(nil)

func (p *postBiz) Create(ctx context.Context, rq *apiv1.CreatePostRequest) (*apiv1.CreatePostResponse, error) {
	var postModel model.Post
	_ = copier.Copy(&postModel, rq)
	// 从ctx中获取到用户ID
	postModel.UserID = contextx.UserID(ctx)

	if err := p.store.Post().Create(ctx, &postModel); err != nil {
		return nil, err
	}

	return &apiv1.CreatePostResponse{PostID: postModel.PostID}, nil
}

func (p *postBiz) Update(ctx context.Context, rq *apiv1.UpdatePostRequest) (*apiv1.UpdatePostResponse, error) {
	whr := where.F("userID", contextx.UserID(ctx), "postID", rq.PostID)
	postModel, err := p.store.Post().Get(ctx, whr)
	if err != nil {
		return nil, err
	}

	if rq.Title != nil {
		postModel.Title = *rq.Title
	}
	if rq.Content != nil {
		postModel.Content = *rq.Content
	}
	if err := p.store.Post().Update(ctx, postModel); err != nil {
		return nil, err
	}

	return &apiv1.UpdatePostResponse{}, nil
}

func (p *postBiz) Delete(ctx context.Context, rq *apiv1.DeletePostRequest) (*apiv1.DeletePostResponse, error) {
	whr := where.F("userID", contextx.UserID(ctx), "postID", rq.PostIDs)
	if err := p.store.Post().Delete(ctx, whr); err != nil {
		return nil, err
	}
	return &apiv1.DeletePostResponse{}, nil
}

func (b *postBiz) Get(ctx context.Context, rq *apiv1.GetPostRequest) (*apiv1.GetPostResponse, error) {
	whr := where.F("userID", contextx.UserID(ctx), "postID", rq.PostID)
	postM, err := b.store.Post().Get(ctx, whr)
	if err != nil {
		return nil, err
	}

	return &apiv1.GetPostResponse{Post: conversion.PostodelToPostV1(postM)}, nil
}

func (b *postBiz) List(ctx context.Context, rq *apiv1.ListPostRequest) (*apiv1.ListPostResponse, error) {
	whr := where.F("userID", contextx.UserID(ctx)).P(int(rq.Offset), int(rq.Limit))
	if rq.Title != nil {
		whr = whr.Q("title like ?", "%"+*rq.Title+"%")
	}

	count, postList, err := b.store.Post().List(ctx, whr)
	if err != nil {
		return nil, err
	}

	posts := make([]*apiv1.Post, 0, len(postList))
	for _, post := range postList {
		converted := conversion.PostodelToPostV1(post)
		posts = append(posts, converted)
	}

	return &apiv1.ListPostResponse{TotalCount: count, Posts: posts}, nil
}
