package services

import (
	"po-backend/models"
	"po-backend/repositories"
)

type FollowService struct {
	Repo *repositories.FollowRepository
}

func NewFollowService(repo *repositories.FollowRepository) *FollowService {
	return &FollowService{Repo: repo}
}

func (s *FollowService) Follow(followerID, followingID uint) (*models.Follow, error) {
	return s.Repo.Follow(followerID, followingID)
}

func (s *FollowService) Unfollow(followerID, followingID uint) error {
	return s.Repo.Unfollow(followerID, followingID)
}

func (s *FollowService) GetFollowingUsers(userID uint) ([]models.User, error) {
	return s.Repo.GetFollowingUsers(userID)
}

func (s *FollowService) GetFollowingIDs(userID uint) ([]uint, error) {
	return s.Repo.GetFollowingIDs(userID)
}

func (s *FollowService) GetFollowerIDs(userID uint) ([]uint, error) {
	return s.Repo.GetFollowerIDs(userID)
}
