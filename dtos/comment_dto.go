package dtos

type CreateCommentRequest struct {
	Content string `json:"content" validate:"required"`
	PostID  uint   `json:"postId" validate:"required"`
}
