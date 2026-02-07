package services

import (
	"po-backend/models"
	"po-backend/repositories"
)

type LikeService struct {
	Repo *repositories.LikeRepository
}

func NewLikeService(repo *repositories.LikeRepository) *LikeService {
	return &LikeService{Repo: repo}
}

func (s *LikeService) LikePost(postID, userID uint) (*models.PostLike, error) {
	return s.Repo.LikePost(postID, userID)
}

func (s *LikeService) UnlikePost(postID, userID uint) error {
	return s.Repo.UnlikePost(postID, userID)
}

func (s *LikeService) LikeComment(commentID, userID uint) (*models.CommentLike, error) {
	return s.Repo.LikeComment(commentID, userID)
}

func (s *LikeService) UnlikeComment(commentID, userID uint) error {
	return s.Repo.UnlikeComment(commentID, userID)
}

func (s *LikeService) GetPostLikers(postID uint) ([]models.PostLike, error) {
	return s.Repo.GetPostLikers(postID)
}

func (s *LikeService) GetCommentLikers(commentID uint) ([]models.CommentLike, error) {
	return s.Repo.GetCommentLikers(commentID)
}
