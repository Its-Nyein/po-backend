package repositories

import (
	"po-backend/models"

	"gorm.io/gorm"
)

type FollowRepository struct {
	DB *gorm.DB
}

func NewFollowRepository(db *gorm.DB) *FollowRepository {
	return &FollowRepository{DB: db}
}

func (r *FollowRepository) Follow(followerID, followingID uint) (*models.Follow, error) {
	follow := &models.Follow{FollowerID: followerID, FollowingID: followingID}
	err := r.DB.Create(follow).Error
	return follow, err
}

func (r *FollowRepository) Unfollow(followerID, followingID uint) error {
	return r.DB.Where("follower_id = ? AND following_id = ?", followerID, followingID).
		Delete(&models.Follow{}).Error
}

func (r *FollowRepository) GetFollowingUsers(userID uint) ([]models.User, error) {
	var users []models.User
	err := r.DB.
		Joins("JOIN follows ON follows.following_id = users.id").
		Where("follows.follower_id = ?", userID).
		Find(&users).Error
	return users, err
}

func (r *FollowRepository) GetFollowingIDs(userID uint) ([]uint, error) {
	var follows []models.Follow
	err := r.DB.Where("follower_id = ?", userID).Find(&follows).Error
	if err != nil {
		return nil, err
	}
	ids := make([]uint, len(follows))
	for i, f := range follows {
		ids[i] = f.FollowingID
	}
	return ids, nil
}
