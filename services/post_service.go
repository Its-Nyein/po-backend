package services

import (
	"po-backend/models"
	"po-backend/repositories"
)

type PostService struct {
	Repo *repositories.PostRepository
}

func NewPostService(repo *repositories.PostRepository) *PostService {
	return &PostService{Repo: repo}
}

func (s *PostService) GetAll() ([]models.Post, error) {
	return s.Repo.GetAll(20)
}

func (s *PostService) GetByID(id uint) (*models.Post, error) {
	return s.Repo.GetByID(id)
}

func (s *PostService) Create(content string, userID uint, quotedPostID *uint) (*models.Post, error) {
	post := &models.Post{
		Content:      content,
		UserID:       userID,
		QuotedPostID: quotedPostID,
	}
	if err := s.Repo.Create(post); err != nil {
		return nil, err
	}
	return s.Repo.GetByID(post.ID)
}

func (s *PostService) Update(id uint, content string) (*models.Post, error) {
	return s.Repo.Update(id, content)
}

func (s *PostService) Delete(id uint) error {
	return s.Repo.Delete(id)
}

func (s *PostService) GetByUserIDs(userIDs []uint) ([]models.Post, error) {
	return s.Repo.GetByUserIDs(userIDs, 20)
}

func (s *PostService) IsOwner(postID, userID uint) bool {
	return s.Repo.IsOwner(postID, userID)
}
