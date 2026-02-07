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
