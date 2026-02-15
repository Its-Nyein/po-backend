package repositories

import (
	"po-backend/models"

	"gorm.io/gorm"
)

type MessageRepository struct {
	DB *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{DB: db}
}

func (r *MessageRepository) GetByConversationID(conversationID uint, cursor uint, limit int) ([]models.Message, error) {
	var messages []models.Message
	query := r.DB.Preload("Sender").
		Where("conversation_id = ?", conversationID).
		Order("created_at DESC").
		Limit(limit)

	if cursor > 0 {
		query = query.Where("id < ?", cursor)
	}

	err := query.Find(&messages).Error
	return messages, err
}

func (r *MessageRepository) Create(message *models.Message) error {
	return r.DB.Create(message).Error
}
