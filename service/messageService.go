package service

import (
	"context"
	"gochatbackend/model"
	"gochatbackend/repository"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MessageService interface {
	CreateMessage(ctx context.Context, message model.Message) (model.MessageResponse, error)
	GetMessageConversationId(ctx context.Context, conversationId primitive.ObjectID, limit int64) ([]model.MessageResponse, error)
}

type messageService struct {
	db                     *mongo.Database
	messageRepository      repository.MessageRepository
	userRepository         repository.UserRepository
	conversationRepository repository.ConversationRepository
}

func NewMessageService(
	db *mongo.Database,
	messageRepository repository.MessageRepository,
	userRepository repository.UserRepository,
	conversationRepository repository.ConversationRepository,
) MessageService {
	return messageService{
		db:                     db,
		messageRepository:      messageRepository,
		userRepository:         userRepository,
		conversationRepository: conversationRepository,
	}
}

func (m messageService) CreateMessage(ctx context.Context, message model.Message) (model.MessageResponse, error) {
	_, err := m.conversationRepository.UserConversation().FindOne(ctx, m.db, map[string]interface{}{
		"userId":         message.UserID,
		"conversationId": message.ConversationID,
	})
	if err != nil {
		log.Error(err)
		return model.MessageResponse{}, err
	}
	result, err := m.messageRepository.Create(ctx, m.db, &message)
	if err != nil {
		log.Error(err)
		return model.MessageResponse{}, err
	}
	user, err := m.userRepository.FindByID(ctx, m.db, message.UserID)
	if err != nil {
		log.Error(err)
		return model.MessageResponse{}, err
	}
	err = m.conversationRepository.Conversation().Update(ctx, m.db, message.ConversationID,
		map[string]interface{}{
			"lastMessageAt": result.CreatedAt,
		},
	)
	if err != nil {
		log.Error(err)
		return model.MessageResponse{}, err
	}
	user.Password = ""
	messageResponse := model.MessageResponse{
		ID:             result.ID,
		CreatedAt:      result.CreatedAt,
		UpdatedAt:      result.UpdatedAt,
		DeletedAt:      result.DeletedAt,
		Content:        result.Content,
		UserID:         result.UserID,
		ConversationID: result.ConversationID,
		UserDetail:     *user,
	}
	return messageResponse, nil
}

func (m messageService) GetMessageConversationId(ctx context.Context, conversationId primitive.ObjectID, limit int64) ([]model.MessageResponse, error) {
	if limit == 0 {
		limit = 20
	}
	result, err := m.messageRepository.GetMessageConversationId(ctx, m.db, conversationId, limit)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return result, nil
}
