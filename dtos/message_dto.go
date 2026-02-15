package dtos

type CreateConversationRequest struct {
	UserID uint `json:"userId" validate:"required"`
}

type SendMessageRequest struct {
	Content string `json:"content" validate:"required,max=2000"`
}
