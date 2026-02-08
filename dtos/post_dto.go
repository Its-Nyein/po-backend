package dtos

type CreatePostRequest struct {
	Content      string `json:"content" validate:"required_without=QuotedPostID"`
	QuotedPostID *uint  `json:"quotedPostId,omitempty"`
}

type UpdatePostRequest struct {
	Content string `json:"content" validate:"required"`
}
