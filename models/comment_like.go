package models

import "time"

type CommentLike struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CommentID uint      `gorm:"not null" json:"commentId"`
	UserID    uint      `gorm:"not null" json:"userId"`
	Comment   Comment   `gorm:"foreignKey:CommentID;constraint:OnDelete:CASCADE" json:"-"`
	User      User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}
