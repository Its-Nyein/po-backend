package dtos

type NotificationResponse struct {
	ID        uint        `json:"id"`
	Type      string      `json:"type"`
	Content   string      `json:"content"`
	UserID    uint        `json:"userId"`
	PostID    uint        `json:"postId"`
	Read      bool        `json:"read"`
	User      interface{} `json:"user"`
	CreatedAt string      `json:"createdAt"`
}
