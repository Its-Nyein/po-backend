package models

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(200);not null" json:"name"`
	Username  string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"username"`
	Bio       string    `gorm:"type:text" json:"bio"`
	Password  string    `gorm:"type:varchar(200);not null" json:"-"`
	Posts     []Post    `gorm:"foreignKey:UserID" json:"posts,omitempty"`
	Comments  []Comment `gorm:"foreignKey:UserID" json:"comments,omitempty"`
	Followers []Follow  `gorm:"foreignKey:FollowingID" json:"followers,omitempty"`
	Following []Follow  `gorm:"foreignKey:FollowerID" json:"following,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
