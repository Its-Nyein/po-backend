package dtos

type CreatePostRequest struct {
	Content string `json:"content" validate:"required"`
}

type UpdatePostRequest struct {
	Content string `json:"content" validate:"required"`
}
