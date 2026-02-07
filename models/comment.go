package models

import "time"

type Comment struct {
	ID        uint          `gorm:"primaryKey" json:"id"`
	Content   string        `gorm:"type:text;not null" json:"content"`
	UserID    uint          `gorm:"not null" json:"userId"`
	PostID    uint          `gorm:"not null" json:"postId"`
	User      User          `gorm:"foreignKey:UserID" json:"user"`
	Post      Post          `gorm:"foreignKey:PostID" json:"-"`
	Likes     []CommentLike `gorm:"foreignKey:CommentID" json:"likes,omitempty"`
	CreatedAt time.Time     `json:"createdAt"`
	UpdatedAt time.Time     `json:"updated_at"`
}
