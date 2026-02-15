package repositories

import (
	"po-backend/models"
	"time"

	"gorm.io/gorm"
)

type ConversationWithPreview struct {
	ID            uint        `json:"id"`
	OtherUser     models.User `json:"otherUser"`
	LastMessage   *string     `json:"lastMessage"`
	LastMessageAt *time.Time  `json:"lastMessageAt"`
	UnreadCount   int64       `json:"unreadCount"`
}

type ConversationRepository struct {
	DB *gorm.DB
}

func NewConversationRepository(db *gorm.DB) *ConversationRepository {
	return &ConversationRepository{DB: db}
}

func (r *ConversationRepository) FindBetweenUsers(userID1, userID2 uint) (*models.Conversation, error) {
	var conversation models.Conversation
	err := r.DB.Where("id IN (?)",
		r.DB.Model(&models.ConversationParticipant{}).
			Select("conversation_id").
			Where("user_id = ?", userID1).
			Where("conversation_id IN (?)",
				r.DB.Model(&models.ConversationParticipant{}).
					Select("conversation_id").
					Where("user_id = ?", userID2),
			),
	).First(&conversation).Error
	if err != nil {
		return nil, err
	}
	return &conversation, nil
}

func (r *ConversationRepository) Create(userID1, userID2 uint) (*models.Conversation, error) {
	var conversation models.Conversation
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&conversation).Error; err != nil {
			return err
		}
		participants := []models.ConversationParticipant{
			{ConversationID: conversation.ID, UserID: userID1},
			{ConversationID: conversation.ID, UserID: userID2},
		}
		return tx.Create(&participants).Error
	})
	if err != nil {
		return nil, err
	}
	return &conversation, nil
}

func (r *ConversationRepository) GetByUserID(userID uint) ([]ConversationWithPreview, error) {
	var participants []models.ConversationParticipant
	err := r.DB.Where("user_id = ?", userID).Find(&participants).Error
	if err != nil {
		return nil, err
	}

	if len(participants) == 0 {
		return []ConversationWithPreview{}, nil
	}

	results := make([]ConversationWithPreview, 0)
	for _, p := range participants {
		var otherParticipant models.ConversationParticipant
		err := r.DB.Preload("User").
			Where("conversation_id = ? AND user_id != ?", p.ConversationID, userID).
			First(&otherParticipant).Error
		if err != nil {
			continue
		}

		preview := ConversationWithPreview{
			ID:        p.ConversationID,
			OtherUser: otherParticipant.User,
		}

		var lastMsg models.Message
		if err := r.DB.Where("conversation_id = ?", p.ConversationID).
			Order("created_at DESC").First(&lastMsg).Error; err == nil {
			preview.LastMessage = &lastMsg.Content
			preview.LastMessageAt = &lastMsg.CreatedAt
		}

		var unreadCount int64
		query := r.DB.Model(&models.Message{}).
			Where("conversation_id = ? AND sender_id != ?", p.ConversationID, userID)
		if p.LastReadAt != nil {
			query = query.Where("created_at > ?", *p.LastReadAt)
		}
		query.Count(&unreadCount)
		preview.UnreadCount = unreadCount

		results = append(results, preview)
	}

	return results, nil
}

func (r *ConversationRepository) GetUnreadTotal(userID uint) (int64, error) {
	var participants []models.ConversationParticipant
	err := r.DB.Where("user_id = ?", userID).Find(&participants).Error
	if err != nil {
		return 0, err
	}

	var total int64
	for _, p := range participants {
		var count int64
		query := r.DB.Model(&models.Message{}).
			Where("conversation_id = ? AND sender_id != ?", p.ConversationID, userID)
		if p.LastReadAt != nil {
			query = query.Where("created_at > ?", *p.LastReadAt)
		}
		query.Count(&count)
		total += count
	}

	return total, nil
}

func (r *ConversationRepository) UpdateLastRead(conversationID, userID uint) error {
	now := time.Now()
	return r.DB.Model(&models.ConversationParticipant{}).
		Where("conversation_id = ? AND user_id = ?", conversationID, userID).
		Update("last_read_at", now).Error
}

func (r *ConversationRepository) IsParticipant(conversationID, userID uint) bool {
	var count int64
	r.DB.Model(&models.ConversationParticipant{}).
		Where("conversation_id = ? AND user_id = ?", conversationID, userID).
		Count(&count)
	return count > 0
}

func (r *ConversationRepository) GetOtherParticipantID(conversationID, userID uint) (uint, error) {
	var participant models.ConversationParticipant
	err := r.DB.Where("conversation_id = ? AND user_id != ?", conversationID, userID).
		First(&participant).Error
	if err != nil {
		return 0, err
	}
	return participant.UserID, nil
}
