package repository

import (
	"context"
	"gochatbackend/model"
	"gochatbackend/pkg/databaseutil"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MessageRepository interface {
	databaseutil.BaseRepository[model.Message]
	GetMessageConversationId(ctx context.Context, db *mongo.Database, conversationId primitive.ObjectID, limit int64) ([]model.MessageResponse, error)
	CreateIndex(ctx context.Context, db *mongo.Database) error
}

type messageRepository struct {
	databaseutil.BaseRepository[model.Message]
}

func NewMessageRepository() MessageRepository {
	return messageRepository{databaseutil.NewBaseRepository[model.Message]("messages")}
}

func (m messageRepository) GetMessageConversationId(ctx context.Context, db *mongo.Database, conversationId primitive.ObjectID, limit int64) ([]model.MessageResponse, error) {
	conversationFilter := bson.M{
		"conversationId": conversationId,
	}
	query := []bson.M{
		{"$match": conversationFilter},
		{"$lookup": bson.M{
			"from":         "users",
			"localField":   "userId",
			"foreignField": "_id",
			"as":           "userDetail",
		}},
		{"$unwind": "$userDetail"},
		{"$sort": bson.M{"createdAt": -1}},
		{"$limit": limit},
	}
	cur, err := db.Collection("messages").Aggregate(ctx, query, options.Aggregate())
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var messageList []model.MessageResponse
	for cur.Next(ctx) {
		var message model.MessageResponse
		err := cur.Decode(&message)
		message.UserDetail.Password = ""
		if err != nil {
			log.Error(err)
			message.Content = "error: invalid message"
			messageList = append(messageList, message)
			continue
		}
		messageList = append(messageList, message)
	}
	return messageList, nil
}

func (m messageRepository) CreateIndex(ctx context.Context, db *mongo.Database) error {
	_, err := db.Collection("messages").Indexes().CreateOne(
		ctx,
		mongo.IndexModel{
			Keys:    bson.M{"conversationId": 1},
			Options: options.Index().SetName("conversationid_index"),
		},
	)
	if err != nil {
		return err
	}
	return nil
}
