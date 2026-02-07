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

func (s *NotificationService) Create(notiType, content string, actorID, userID uint, postID *uint) (*models.Notification, error) {
	noti := &models.Notification{
		Type:    notiType,
		Content: content,
		ActorID: actorID,
		UserID:  userID,
		PostID:  postID,
	}
	if err := s.Repo.Create(noti); err != nil {
		return nil, err
	}
	return noti, nil
}

func (s *NotificationService) GetByUserID(userID uint) ([]models.Notification, error) {
	return s.Repo.GetByUserID(userID)
}

func (s *NotificationService) MarkAllRead(userID uint) error {
	return s.Repo.MarkAllRead(userID)
}

func (s *NotificationService) MarkOneRead(id uint) error {
	return s.Repo.MarkOneRead(id)
}
