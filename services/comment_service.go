package services

import (
	"po-backend/models"
	"po-backend/repositories"
)

type CommentService struct {
	Repo *repositories.CommentRepository
}

func NewCommentService(repo *repositories.CommentRepository) *CommentService {
	return &CommentService{Repo: repo}
}

func (s *CommentService) Create(content string, userID, postID uint) (*models.Comment, error) {
	comment := &models.Comment{
		Content: content,
		UserID:  userID,
		PostID:  postID,
	}
	if err := s.Repo.Create(comment); err != nil {
		return nil, err
	}
	return s.Repo.GetByID(comment.ID)
}

func (s *CommentService) Update(id uint, content string) (*models.Comment, error) {
	return s.Repo.Update(id, content)
}

func (s *CommentService) Delete(id uint) error {
	return s.Repo.Delete(id)
}

func (s *CommentService) IsOwner(commentID, userID uint) bool {
	return s.Repo.IsOwner(commentID, userID)
}

func (s *CommentService) GetByID(id uint) (*models.Comment, error) {
	return s.Repo.GetByID(id)
}
