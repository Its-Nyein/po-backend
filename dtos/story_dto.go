package dtos

type CreateStoryRequest struct {
	Content string `json:"content" validate:"required,max=500"`
	Privacy string `json:"privacy" validate:"required,oneof=public friends private"`
}
