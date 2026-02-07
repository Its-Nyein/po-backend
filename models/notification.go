package models

import "time"

type Notification struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Type      string    `gorm:"type:varchar(50);not null" json:"type"`
	Content   string    `gorm:"type:text" json:"content"`
	UserID    uint      `gorm:"not null" json:"userId"`
	ActorID   uint      `gorm:"not null" json:"actorId"`
	PostID    *uint     `gorm:"default:null" json:"postId"`
	Read      bool      `gorm:"default:false" json:"read"`
	Actor     User      `gorm:"foreignKey:ActorID;constraint:OnDelete:CASCADE" json:"actor,omitempty"`
	User      User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	Post      *Post     `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE" json:"post,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}
