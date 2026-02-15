package models

import "time"

type ConversationParticipant struct {
	ID             uint         `gorm:"primaryKey" json:"id"`
	ConversationID uint         `gorm:"not null;uniqueIndex:idx_conv_user" json:"conversationId"`
	UserID         uint         `gorm:"not null;uniqueIndex:idx_conv_user" json:"userId"`
	LastReadAt     *time.Time   `json:"lastReadAt"`
	Conversation   Conversation `gorm:"foreignKey:ConversationID;constraint:OnDelete:CASCADE" json:"-"`
	User           User         `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	CreatedAt      time.Time    `json:"createdAt"`
}
