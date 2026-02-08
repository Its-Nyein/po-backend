package repositories

import (
	"po-backend/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) GetAll(limit int) ([]models.User, error) {
	var users []models.User
	err := r.DB.Preload("Followers").Preload("Following").
		Order("created_at DESC").Limit(limit).Find(&users).Error
	return users, err
}

func (r *UserRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	err := r.DB.Preload("Followers").Preload("Following").
		Preload("Posts", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC")
		}).
		Preload("Posts.User").
		Preload("Posts.Comments.User").
		Preload("Posts.Comments.Likes").
		Preload("Posts.Likes").
		Preload("Posts.QuotedPost.User").
		First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Create(user *models.User) error {
	return r.DB.Create(user).Error
}

func (r *UserRepository) Search(query string, limit int) ([]models.User, error) {
	var users []models.User
	err := r.DB.Preload("Followers").Preload("Following").
		Where("name ILIKE ?", "%"+query+"%").
		Limit(limit).Find(&users).Error
	return users, err
}

func (r *UserRepository) Update(user *models.User) error {
	return r.DB.Save(user).Error
}

func (r *UserRepository) Delete(id uint) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("DELETE FROM comment_likes WHERE comment_id IN (SELECT id FROM comments WHERE user_id = ?)", id).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM comment_likes WHERE comment_id IN (SELECT id FROM comments WHERE post_id IN (SELECT id FROM posts WHERE user_id = ?))", id).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM post_likes WHERE post_id IN (SELECT id FROM posts WHERE user_id = ?)", id).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM comments WHERE post_id IN (SELECT id FROM posts WHERE user_id = ?)", id).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM bookmarks WHERE post_id IN (SELECT id FROM posts WHERE user_id = ?)", id).Error; err != nil {
			return err
		}
		tables := []string{"post_likes", "comment_likes", "comments", "bookmarks", "notifications"}
		for _, table := range tables {
			if err := tx.Exec("DELETE FROM "+table+" WHERE user_id = ?", id).Error; err != nil {
				return err
			}
		}
		if err := tx.Exec("DELETE FROM notifications WHERE actor_id = ?", id).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM follows WHERE follower_id = ? OR following_id = ?", id, id).Error; err != nil {
			return err
		}
		if err := tx.Exec("UPDATE posts SET quoted_post_id = NULL WHERE quoted_post_id IN (SELECT id FROM posts WHERE user_id = ?)", id).Error; err != nil {
			return err
		}
		if err := tx.Where("user_id = ?", id).Delete(&models.Post{}).Error; err != nil {
			return err
		}
		if err := tx.Delete(&models.User{}, id).Error; err != nil {
			return err
		}
		return nil
	})
}
