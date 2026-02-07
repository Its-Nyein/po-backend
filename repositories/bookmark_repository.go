package repositories

import (
	"po-backend/models"

	"gorm.io/gorm"
)

type BookmarkRepository struct {
	DB *gorm.DB
}

func NewBookmarkRepository(db *gorm.DB) *BookmarkRepository {
	return &BookmarkRepository{DB: db}
}

func (r *BookmarkRepository) Create(postID, userID uint) (*models.Bookmark, error) {
	bookmark := &models.Bookmark{
		PostID: postID,
		UserID: userID,
	}
	if err := r.DB.Create(bookmark).Error; err != nil {
		return nil, err
	}
	return bookmark, nil
}

func (r *BookmarkRepository) Delete(postID, userID uint) error {
	return r.DB.Where("post_id = ? AND user_id = ?", postID, userID).Delete(&models.Bookmark{}).Error
}

func (r *BookmarkRepository) GetByUserID(userID uint) ([]models.Bookmark, error) {
	var bookmarks []models.Bookmark
	err := r.DB.Preload("Post.User").
		Preload("Post.Comments.User").
		Preload("Post.Comments.Likes").
		Preload("Post.Likes").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&bookmarks).Error
	return bookmarks, err
}

func (r *BookmarkRepository) Exists(postID, userID uint) bool {
	var count int64
	r.DB.Model(&models.Bookmark{}).Where("post_id = ? AND user_id = ?", postID, userID).Count(&count)
	return count > 0
}
