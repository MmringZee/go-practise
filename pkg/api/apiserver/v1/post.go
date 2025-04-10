// Post API 定义，包含博客文章的请求和响应消息

package v1

import (
	"time"
)

// 博客文章
type Post struct {
	// 博文 ID
	PostID string `json:"postID"`
	// 用户 ID
	UserID string `json:"userID"`
	// 博客标题
	Title string `json:"title"`
	// 博客内容
	Content string `json:"content"`
	// 博客创建时间
	CreatedAt time.Time `json:"createdAt"`
	// 博客最后更新时间
	UpdatedAt time.Time `json:"updatedAt"`
}

// 创建文章请求
type CreatePostRequest struct {
	// 博客标题
	Title string `json:"title"`
	// 博客内容
	Content string `json:"content"`
}

// 创建文章响应
type CreatePostResponse struct {
	// 创建的文章 ID
	PostID string `json:"postID"`
}

// 更新文章请求
type UpdatePostRequest struct {
	// 要更新的文章 ID，对应 {postID}
	PostID string `json:"postID" uri:"postID"`
	// 更新后的博客标题
	Title *string `json:"title"`
	// 更新后的博客内容
	Content *string `json:"content"`
}

// 更新文章响应
type UpdatePostResponse struct {
}

// 删除文章请求
type DeletePostRequest struct {
	// 要删除的文章 ID 列表
	PostIDs []string `json:"postIDs"`
}

// 删除文章响应
type DeletePostResponse struct {
}

// 获取文章请求
type GetPostRequest struct {
	// 要获取的文章 ID
	PostID string `json:"postID" uri:"postID"`
}

// 获取文章响应
type GetPostResponse struct {
	// 返回的文章信息
	Post *Post `json:"post"`
}

// 获取文章列表请求
type ListPostRequest struct {
	// 偏移量
	Offset int64 `json:"offset"`
	// 每页数量
	Limit int64 `json:"limit"`
	// 可选的标题过滤
	Title *string `json:"title"`
}

// 获取文章列表响应
type ListPostResponse struct {
	// 总文章数
	TotalCount int64 `json:"totalCount"`
	// 文章列表
	Posts []*Post `json:"posts"`
}
