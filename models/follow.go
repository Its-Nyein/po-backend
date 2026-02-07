package models

import "time"

type Follow struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	FollowerID  uint      `gorm:"not null" json:"followerId"`
	FollowingID uint      `gorm:"not null" json:"followingId"`
	Follower    User      `gorm:"foreignKey:FollowerID;constraint:OnDelete:CASCADE" json:"follower,omitempty"`
	Following   User      `gorm:"foreignKey:FollowingID;constraint:OnDelete:CASCADE" json:"following,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
}
