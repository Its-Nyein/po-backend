package repositories

import (
	"po-backend/models"

	"gorm.io/gorm"
)

type CommentRepository struct {
	DB *gorm.DB
}

func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{DB: db}
}

func (r *CommentRepository) Create(comment *models.Comment) error {
	return r.DB.Create(comment).Error
}

func (r *CommentRepository) GetByID(id uint) (*models.Comment, error) {
	var comment models.Comment
	err := r.DB.Preload("User").Preload("Likes").First(&comment, id).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *CommentRepository) Delete(id uint) error {
	return r.DB.Delete(&models.Comment{}, id).Error
}

func (r *CommentRepository) Update(id uint, content string) (*models.Comment, error) {
	if err := r.DB.Model(&models.Comment{}).Where("id = ?", id).Update("content", content).Error; err != nil {
		return nil, err
	}
	return r.GetByID(id)
}

func (r *CommentRepository) IsOwner(commentID, userID uint) bool {
	var comment models.Comment
	err := r.DB.Where("id = ? AND user_id = ?", commentID, userID).First(&comment).Error
	return err == nil
}
