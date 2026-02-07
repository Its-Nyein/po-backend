package services

import (
	"po-backend/models"
	"po-backend/repositories"
)

type NotificationService struct {
	Repo *repositories.NotificationRepository
}

func NewNotificationService(repo *repositories.NotificationRepository) *NotificationService {
	return &NotificationService{Repo: repo}
}

func (s *NotificationService) Create(notiType, content string, userID, postID uint) (*models.Notification, error) {
	noti := &models.Notification{
		Type:    notiType,
		Content: content,
		UserID:  userID,
		PostID:  postID,
	}
	if err := s.Repo.Create(noti); err != nil {
		return nil, err
	}
	return noti, nil
}

func (s *NotificationService) GetByPostOwner(userID uint) ([]models.Notification, error) {
	return s.Repo.GetByPostOwner(userID)
}

func (s *NotificationService) MarkAllRead(userID uint) error {
	return s.Repo.MarkAllRead(userID)
}

func (s *NotificationService) MarkOneRead(id uint) error {
	return s.Repo.MarkOneRead(id)
}
