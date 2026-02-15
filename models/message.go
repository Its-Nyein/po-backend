package models

import "time"

type Message struct {
	ID             uint         `gorm:"primaryKey" json:"id"`
	ConversationID uint         `gorm:"not null;index" json:"conversationId"`
	SenderID       uint         `gorm:"not null" json:"senderId"`
	Content        string       `gorm:"type:text;not null" json:"content"`
	Conversation   Conversation `gorm:"foreignKey:ConversationID;constraint:OnDelete:CASCADE" json:"-"`
	Sender         User         `gorm:"foreignKey:SenderID;constraint:OnDelete:CASCADE" json:"sender,omitempty"`
	CreatedAt      time.Time    `json:"createdAt"`
}
