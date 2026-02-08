package services

import (
	"po-backend/models"
	"po-backend/repositories"
	"po-backend/utilities"
)

type HashtagService struct {
	Repo *repositories.HashtagRepository
}

func NewHashtagService(repo *repositories.HashtagRepository) *HashtagService {
	return &HashtagService{Repo: repo}
}

func (s *HashtagService) SyncPostHashtags(postID uint, content string) error {
	tags := utilities.ParseHashtags(content)
	if len(tags) == 0 {
		return s.Repo.SyncPostHashtags(postID, []string{})
	}
	return s.Repo.SyncPostHashtags(postID, tags)
}

func (s *HashtagService) GetPostsByHashtag(tag string) ([]models.Post, error) {
	return s.Repo.GetPostsByHashtag(tag, 50)
}

func (s *HashtagService) GetTrending(limit int) ([]map[string]interface{}, error) {
	return s.Repo.GetTrending(limit)
}
