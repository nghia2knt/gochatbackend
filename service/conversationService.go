package service

import (
	"context"
	"gochatbackend/model"
	"gochatbackend/repository"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ConversationService interface {
	CreateConversation(ctx context.Context, name string, userId primitive.ObjectID, userIdList []primitive.ObjectID) (model.Conversation, error)
	GetConversationUserId(ctx context.Context, userId primitive.ObjectID, limit int64) ([]model.Conversation, error)
	GetConversationById(ctx context.Context, userId primitive.ObjectID, conversationId primitive.ObjectID) (model.ConversationResponse, error)
}

type conversationService struct {
	db                     *mongo.Database
	conversationRepository repository.ConversationRepository
}

func NewConversationService(
	db *mongo.Database,
	conversationRepository repository.ConversationRepository,
) ConversationService {
	return conversationService{
		db:                     db,
		conversationRepository: conversationRepository,
	}
}

func (c conversationService) CreateConversation(ctx context.Context, name string, userId primitive.ObjectID, userIdList []primitive.ObjectID) (model.Conversation, error) {
	conversation, err := c.conversationRepository.CreateConversationWithOwnership(ctx, c.db, name, userId, userIdList)
	if err != nil {
		log.Error(err)
		return model.Conversation{}, err
	}
	return conversation, nil
}

func (c conversationService) GetConversationUserId(ctx context.Context, userId primitive.ObjectID, limit int64) ([]model.Conversation, error) {
	if limit == 0 {
		limit = 20
	}
	result, err := c.conversationRepository.GetConversationByUserId(ctx, c.db, userId, limit)
	if err != nil {
		log.Error(err)
		return []model.Conversation{}, err
	}
	return result, nil
}

func (c conversationService) GetConversationById(ctx context.Context, userId primitive.ObjectID, conversationId primitive.ObjectID) (model.ConversationResponse, error) {
	_, err := c.conversationRepository.UserConversation().FindOne(ctx, c.db, map[string]interface{}{
		"userId":         userId,
		"conversationId": conversationId,
	})
	if err != nil {
		log.Error(err)
		return model.ConversationResponse{}, err
	}
	result, err := c.conversationRepository.Conversation().FindByID(ctx, c.db, conversationId)
	if err != nil {
		log.Error(err)
		return model.ConversationResponse{}, err
	}
	listUserConversation, err := c.conversationRepository.UserConversation().Find(ctx, c.db, map[string]interface{}{
		"conversationId": conversationId,
	})
	if err != nil {
		log.Error(err)
		return model.ConversationResponse{}, err
	}
	var listUser []primitive.ObjectID
	for _, value := range listUserConversation {
		listUser = append(listUser, value.UserID)
	}
	conversationResponse := model.ConversationResponse{
		ID:            result.ID,
		CreatedAt:     result.CreatedAt,
		UpdatedAt:     result.UpdatedAt,
		DeletedAt:     result.DeletedAt,
		Name:          result.Name,
		LastMessageAt: result.LastMessageAt,
		Members:       listUser,
	}
	return conversationResponse, nil
}
