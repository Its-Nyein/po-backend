package repositories

import (
	"po-backend/models"

	"gorm.io/gorm"
)

type LikeRepository struct {
	DB *gorm.DB
}

func NewLikeRepository(db *gorm.DB) *LikeRepository {
	return &LikeRepository{DB: db}
}

func (r *LikeRepository) LikePost(postID, userID uint) (*models.PostLike, error) {
	like := &models.PostLike{PostID: postID, UserID: userID}
	err := r.DB.Create(like).Error
	return like, err
}

func (r *LikeRepository) UnlikePost(postID, userID uint) error {
	return r.DB.Where("post_id = ? AND user_id = ?", postID, userID).
		Delete(&models.PostLike{}).Error
}

func (r *LikeRepository) LikeComment(commentID, userID uint) (*models.CommentLike, error) {
	like := &models.CommentLike{CommentID: commentID, UserID: userID}
	err := r.DB.Create(like).Error
	return like, err
}

func (r *LikeRepository) UnlikeComment(commentID, userID uint) error {
	return r.DB.Where("comment_id = ? AND user_id = ?", commentID, userID).
		Delete(&models.CommentLike{}).Error
}

func (r *LikeRepository) GetPostLikers(postID uint) ([]models.PostLike, error) {
	var likes []models.PostLike
	err := r.DB.Preload("User").Where("post_id = ?", postID).Find(&likes).Error
	return likes, err
}

func (r *LikeRepository) GetCommentLikers(commentID uint) ([]models.CommentLike, error) {
	var likes []models.CommentLike
	err := r.DB.Preload("User").Where("comment_id = ?", commentID).Find(&likes).Error
	return likes, err
}
