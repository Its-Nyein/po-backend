package models

import "time"

type Notification struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Type      string    `gorm:"type:varchar(50);not null" json:"type"`
	Content   string    `gorm:"type:text" json:"content"`
	UserID    uint      `gorm:"not null" json:"userId"`
	PostID    uint      `gorm:"not null" json:"postId"`
	Read      bool      `gorm:"default:false" json:"read"`
	User      User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	Post      Post      `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE" json:"post,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}
