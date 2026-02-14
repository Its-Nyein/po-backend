package models

import "time"

type Story struct {
	ID        uint        `gorm:"primaryKey" json:"id"`
	Content   string      `gorm:"type:text;not null" json:"content"`
	Privacy   string      `gorm:"type:varchar(20);not null;default:'public'" json:"privacy"`
	UserID    uint        `gorm:"not null" json:"userId"`
	User      User        `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user"`
	Views     []StoryView `gorm:"foreignKey:StoryID" json:"views,omitempty"`
	ExpiresAt time.Time   `gorm:"not null;index" json:"expiresAt"`
	CreatedAt time.Time   `json:"createdAt"`
}
