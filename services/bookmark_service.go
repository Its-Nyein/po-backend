package services

import (
	"po-backend/models"
	"po-backend/repositories"
)

type BookmarkService struct {
	Repo *repositories.BookmarkRepository
}

func NewBookmarkService(repo *repositories.BookmarkRepository) *BookmarkService {
	return &BookmarkService{Repo: repo}
}

func (s *BookmarkService) Create(postID, userID uint) (*models.Bookmark, error) {
	return s.Repo.Create(postID, userID)
}

func (s *BookmarkService) Delete(postID, userID uint) error {
	return s.Repo.Delete(postID, userID)
}

func (s *BookmarkService) GetByUserID(userID uint) ([]models.Bookmark, error) {
	return s.Repo.GetByUserID(userID)
}

func (s *BookmarkService) Exists(postID, userID uint) bool {
	return s.Repo.Exists(postID, userID)
}
