package repositories

import (
	"po-backend/models"

	"gorm.io/gorm"
)

type NotificationRepository struct {
	DB *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{DB: db}
}

func (r *NotificationRepository) Create(notification *models.Notification) error {
	return r.DB.Create(notification).Error
}

func (r *NotificationRepository) GetByPostOwner(userID uint) ([]models.Notification, error) {
	var notifications []models.Notification
	err := r.DB.Preload("User").
		Joins("JOIN posts ON posts.id = notifications.post_id").
		Where("posts.user_id = ?", userID).
		Order("notifications.created_at DESC").
		Find(&notifications).Error
	return notifications, err
}

func (r *NotificationRepository) MarkAllRead(userID uint) error {
	return r.DB.Model(&models.Notification{}).
		Joins("JOIN posts ON posts.id = notifications.post_id").
		Where("posts.user_id = ?", userID).
		Update("read", true).Error
}

func (r *NotificationRepository) MarkOneRead(id uint) error {
	return r.DB.Model(&models.Notification{}).Where("id = ?", id).
		Update("read", true).Error
}
