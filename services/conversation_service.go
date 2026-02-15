package services

import (
	"errors"
	"po-backend/models"
	"po-backend/repositories"
)

type ConversationService struct {
	ConvRepo   *repositories.ConversationRepository
	MsgRepo    *repositories.MessageRepository
	FollowRepo *repositories.FollowRepository
}

func NewConversationService(
	convRepo *repositories.ConversationRepository,
	msgRepo *repositories.MessageRepository,
	followRepo *repositories.FollowRepository,
) *ConversationService {
	return &ConversationService{
		ConvRepo:   convRepo,
		MsgRepo:    msgRepo,
		FollowRepo: followRepo,
	}
}

func (s *ConversationService) AreMutualFollowers(userID1, userID2 uint) (bool, error) {
	followingIDs, err := s.FollowRepo.GetFollowingIDs(userID1)
	if err != nil {
		return false, err
	}

	followerIDs, err := s.FollowRepo.GetFollowerIDs(userID1)
	if err != nil {
		return false, err
	}

	followsOther := false
	for _, id := range followingIDs {
		if id == userID2 {
			followsOther = true
			break
		}
	}

	otherFollowsBack := false
	for _, id := range followerIDs {
		if id == userID2 {
			otherFollowsBack = true
			break
		}
	}

	return followsOther && otherFollowsBack, nil
}

func (s *ConversationService) GetOrCreateConversation(currentUserID, otherUserID uint) (*models.Conversation, error) {
	if currentUserID == otherUserID {
		return nil, errors.New("cannot message yourself")
	}

	mutual, err := s.AreMutualFollowers(currentUserID, otherUserID)
	if err != nil {
		return nil, err
	}
	if !mutual {
		return nil, errors.New("you can only message mutual followers")
	}

	conv, err := s.ConvRepo.FindBetweenUsers(currentUserID, otherUserID)
	if err == nil {
		return conv, nil
	}

	return s.ConvRepo.Create(currentUserID, otherUserID)
}

func (s *ConversationService) GetConversations(userID uint) ([]repositories.ConversationWithPreview, error) {
	return s.ConvRepo.GetByUserID(userID)
}

func (s *ConversationService) GetMessages(conversationID, userID uint, cursor uint) ([]models.Message, error) {
	if !s.ConvRepo.IsParticipant(conversationID, userID) {
		return nil, errors.New("not a participant in this conversation")
	}
	return s.MsgRepo.GetByConversationID(conversationID, cursor, 30)
}

func (s *ConversationService) SendMessage(conversationID, senderID uint, content string) (*models.Message, error) {
	if !s.ConvRepo.IsParticipant(conversationID, senderID) {
		return nil, errors.New("not a participant in this conversation")
	}

	otherUserID, err := s.ConvRepo.GetOtherParticipantID(conversationID, senderID)
	if err != nil {
		return nil, err
	}

	mutual, err := s.AreMutualFollowers(senderID, otherUserID)
	if err != nil {
		return nil, err
	}
	if !mutual {
		return nil, errors.New("you can only message mutual followers")
	}

	msg := &models.Message{
		ConversationID: conversationID,
		SenderID:       senderID,
		Content:        content,
	}
	if err := s.MsgRepo.Create(msg); err != nil {
		return nil, err
	}

	messages, err := s.MsgRepo.GetByConversationID(conversationID, 0, 1)
	if err != nil || len(messages) == 0 {
		return msg, nil
	}
	return &messages[0], nil
}

func (s *ConversationService) MarkRead(conversationID, userID uint) error {
	if !s.ConvRepo.IsParticipant(conversationID, userID) {
		return errors.New("not a participant in this conversation")
	}
	return s.ConvRepo.UpdateLastRead(conversationID, userID)
}

func (s *ConversationService) GetUnreadTotal(userID uint) (int64, error) {
	return s.ConvRepo.GetUnreadTotal(userID)
}

func (s *ConversationService) GetOtherParticipantID(conversationID, userID uint) (uint, error) {
	return s.ConvRepo.GetOtherParticipantID(conversationID, userID)
}

func (s *ConversationService) IsParticipant(conversationID, userID uint) bool {
	return s.ConvRepo.IsParticipant(conversationID, userID)
}
