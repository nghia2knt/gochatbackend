package repository

import (
	"context"
	"gochatbackend/model"
	"gochatbackend/pkg/databaseutil"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ConversationRepository interface {
	Conversation() databaseutil.BaseRepository[model.Conversation]
	UserConversation() databaseutil.BaseRepository[model.UserConversation]
	CreateConversationWithOwnership(ctx context.Context, db *mongo.Database, name string, userId primitive.ObjectID, userIdList []primitive.ObjectID) (model.Conversation, error)
	GetConversationByUserId(ctx context.Context, db *mongo.Database, userId primitive.ObjectID, limit int64) ([]model.Conversation, error)
	CreateIndex(ctx context.Context, db *mongo.Database) error
}

type conversationRepository struct {
	conversation     databaseutil.BaseRepository[model.Conversation]
	userConversation databaseutil.BaseRepository[model.UserConversation]
}

func (c conversationRepository) Conversation() databaseutil.BaseRepository[model.Conversation] {
	return c.conversation
}

func (c conversationRepository) UserConversation() databaseutil.BaseRepository[model.UserConversation] {
	return c.userConversation
}

func NewConversationRepository() ConversationRepository {
	return conversationRepository{
		conversation:     databaseutil.NewBaseRepository[model.Conversation]("conversations"),
		userConversation: databaseutil.NewBaseRepository[model.UserConversation]("userConversations"),
	}
}

func (c conversationRepository) CreateConversationWithOwnership(ctx context.Context, db *mongo.Database, name string, userId primitive.ObjectID, userIdList []primitive.ObjectID) (model.Conversation, error) {
	conver := model.Conversation{CreatedAt: time.Now(), Name: name, LastMessageAt: time.Now()}
	conversationInsert, err := db.Collection("conversations").InsertOne(ctx, conver)
	if err != nil {
		return model.Conversation{}, err
	}
	conver.ID = conversationInsert.InsertedID.(primitive.ObjectID)
	var userConversations []interface{}
	mapUserId := map[primitive.ObjectID]bool{}
	userIdList = append(userIdList, userId)
	for _, memberId := range userIdList {
		if mapUserId[memberId] {
			continue
		}
		uC := model.UserConversation{
			CreatedAt:      time.Now(),
			ConversationID: conversationInsert.InsertedID.(primitive.ObjectID),
			UserID:         memberId,
		}
		if memberId == userId {
			uC.Role = model.AdminRole
		}
		mapUserId[memberId] = true
		userConversations = append(userConversations, uC)
	}
	_, err = db.Collection("userConversations").InsertMany(ctx, userConversations)
	if err != nil {
		return model.Conversation{}, err
	}
	return conver, nil
}

func (c conversationRepository) GetConversationByUserId(ctx context.Context, db *mongo.Database, userId primitive.ObjectID, limit int64) ([]model.Conversation, error) {
	userConversations := db.Collection("userConversations")
	pipeline := []bson.M{
		{"$match": bson.M{"userId": userId}},
		{"$lookup": bson.M{
			"from":         "conversations",
			"localField":   "conversationId",
			"foreignField": "_id",
			"as":           "conversation",
		}},
		{"$unwind": "$conversation"},
		{"$project": bson.M{"conversation._id": 1, "conversation.name": 1, "conversation.createdAt": 1, "conversation.lastMessageAt": 1}},
		{"$sort": bson.M{"conversation.lastMessageAt": -1}},
		{"$limit": limit},
	}
	cur, err := userConversations.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var result []model.Conversation
	for cur.Next(ctx) {
		var convData struct {
			Conversation model.Conversation `bson:"conversation"`
		}
		if err := cur.Decode(&convData); err != nil {
			return nil, err
		}
		result = append(result, convData.Conversation)
	}
	return result, nil
}

func (c conversationRepository) CreateIndex(ctx context.Context, db *mongo.Database) error {
	_, err := db.Collection("userConversations").Indexes().CreateOne(
		ctx,
		mongo.IndexModel{
			Keys:    bson.D{{Key: "userId", Value: 1}, {Key: "conversationId", Value: 1}},
			Options: options.Index().SetName("userconversation_index").SetUnique(true),
		},
	)
	if err != nil {
		return err
	}
	return nil
}
