package models

type Hashtag struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"type:varchar(100);uniqueIndex;not null" json:"name"`
}

type PostHashtag struct {
	PostID    uint    `gorm:"primaryKey" json:"postId"`
	HashtagID uint    `gorm:"primaryKey" json:"hashtagId"`
	Post      Post    `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE" json:"-"`
	Hashtag   Hashtag `gorm:"foreignKey:HashtagID;constraint:OnDelete:CASCADE" json:"-"`
}
