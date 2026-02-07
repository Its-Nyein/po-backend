package repositories

import (
	"po-backend/models"

	"gorm.io/gorm"
)

type PostRepository struct {
	DB *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{DB: db}
}

func (r *PostRepository) GetAll(limit int) ([]models.Post, error) {
	var posts []models.Post
	err := r.DB.Preload("User").
		Preload("Comments.User").
		Preload("Comments.Likes").
		Preload("Likes").
		Order("created_at DESC").Limit(limit).Find(&posts).Error
	return posts, err
}

func (r *PostRepository) GetByID(id uint) (*models.Post, error) {
	var post models.Post
	err := r.DB.Preload("User").
		Preload("Comments.User").
		Preload("Comments.Likes").
		Preload("Likes").
		First(&post, id).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *PostRepository) Create(post *models.Post) error {
	return r.DB.Create(post).Error
}

func (r *PostRepository) Delete(id uint) error {
	return r.DB.Delete(&models.Post{}, id).Error
}

func (r *PostRepository) GetByUserIDs(userIDs []uint, limit int) ([]models.Post, error) {
	var posts []models.Post
	err := r.DB.Preload("User").
		Preload("Comments.User").
		Preload("Comments.Likes").
		Preload("Likes").
		Where("user_id IN ?", userIDs).
		Order("created_at DESC").Limit(limit).Find(&posts).Error
	return posts, err
}

func (r *PostRepository) Update(id uint, content string) (*models.Post, error) {
	if err := r.DB.Model(&models.Post{}).Where("id = ?", id).Update("content", content).Error; err != nil {
		return nil, err
	}
	return r.GetByID(id)
}

func (r *PostRepository) IsOwner(postID, userID uint) bool {
	var post models.Post
	err := r.DB.Where("id = ? AND user_id = ?", postID, userID).First(&post).Error
	return err == nil
}
