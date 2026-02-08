package models

import "time"

type Post struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	Content      string     `gorm:"type:text;not null" json:"content"`
	UserID       uint       `gorm:"not null" json:"userId"`
	User         User       `gorm:"foreignKey:UserID" json:"user"`
	Comments     []Comment  `gorm:"foreignKey:PostID" json:"comments,omitempty"`
	Likes        []PostLike `gorm:"foreignKey:PostID" json:"likes,omitempty"`
	QuotedPostID *uint      `json:"quotedPostId,omitempty"`
	QuotedPost   *Post      `gorm:"foreignKey:QuotedPostID" json:"quotedPost,omitempty"`
	QuoteCount   int        `gorm:"-" json:"quoteCount"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
