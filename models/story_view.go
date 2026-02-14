package models

import "time"

type StoryView struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	StoryID   uint      `gorm:"not null;uniqueIndex:idx_story_viewer" json:"storyId"`
	ViewerID  uint      `gorm:"not null;uniqueIndex:idx_story_viewer" json:"viewerId"`
	Story     Story     `gorm:"foreignKey:StoryID;constraint:OnDelete:CASCADE" json:"-"`
	Viewer    User      `gorm:"foreignKey:ViewerID;constraint:OnDelete:CASCADE" json:"viewer,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}
