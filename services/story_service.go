package services

import (
	"po-backend/models"
	"po-backend/repositories"
	"time"
)

type StoryService struct {
	Repo       *repositories.StoryRepository
	FollowRepo *repositories.FollowRepository
}

func NewStoryService(repo *repositories.StoryRepository, followRepo *repositories.FollowRepository) *StoryService {
	return &StoryService{Repo: repo, FollowRepo: followRepo}
}

func (s *StoryService) Create(content, privacy string, userID uint) (*models.Story, error) {
	story := &models.Story{
		Content:   content,
		Privacy:   privacy,
		UserID:    userID,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	if err := s.Repo.Create(story); err != nil {
		return nil, err
	}
	return s.Repo.GetByID(story.ID)
}

func (s *StoryService) Delete(id uint) error {
	return s.Repo.Delete(id)
}

func (s *StoryService) GetFeedStories(viewerID uint) ([]models.Story, error) {
	followingIDs, err := s.FollowRepo.GetFollowingIDs(viewerID)
	if err != nil {
		return nil, err
	}

	mutualIDs, err := s.getMutualIDs(viewerID, followingIDs)
	if err != nil {
		return nil, err
	}

	return s.Repo.GetFeedStories(viewerID, followingIDs, mutualIDs)
}

func (s *StoryService) GetUserStories(targetUserID, viewerID uint) ([]models.Story, error) {
	followingIDs, err := s.FollowRepo.GetFollowingIDs(viewerID)
	if err != nil {
		return nil, err
	}

	mutualIDs, err := s.getMutualIDs(viewerID, followingIDs)
	if err != nil {
		return nil, err
	}

	return s.Repo.GetUserStories(targetUserID, viewerID, mutualIDs)
}

func (s *StoryService) CreateView(storyID, viewerID uint) (*models.StoryView, error) {
	return s.Repo.CreateView(storyID, viewerID)
}

func (s *StoryService) GetViewers(storyID uint) ([]models.StoryView, error) {
	return s.Repo.GetViewers(storyID)
}

func (s *StoryService) IsOwner(storyID, userID uint) bool {
	return s.Repo.IsOwner(storyID, userID)
}

func (s *StoryService) getMutualIDs(viewerID uint, followingIDs []uint) ([]uint, error) {
	if len(followingIDs) == 0 {
		return nil, nil
	}

	followerIDs, err := s.FollowRepo.GetFollowerIDs(viewerID)
	if err != nil {
		return nil, err
	}

	followerSet := make(map[uint]bool, len(followerIDs))
	for _, id := range followerIDs {
		followerSet[id] = true
	}

	var mutualIDs []uint
	for _, id := range followingIDs {
		if followerSet[id] {
			mutualIDs = append(mutualIDs, id)
		}
	}
	return mutualIDs, nil
}
