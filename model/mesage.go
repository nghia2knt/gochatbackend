package model

import (
	"gochatbackend/pkg/databaseutil"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CreatedAt      time.Time          `bson:"createdAt"  json:"createdAt"`
	UpdatedAt      time.Time          `bson:"updatedAt"  json:"updatedAt"`
	DeletedAt      time.Time          `bson:"deletedAt"  json:"deletedAt"`
	Content        string             `bson:"content"  json:"content"`
	UserID         primitive.ObjectID `bson:"userId"  json:"userId"`
	ConversationID primitive.ObjectID `bson:"conversationId" json:"conversationId"`
}

func (u Message) GetBaseModel(model interface{}) databaseutil.BaseModel {
	return databaseutil.BaseModel{
		ID:        &model.(*Message).ID,
		CreatedAt: &model.(*Message).CreatedAt,
		UpdatedAt: &model.(*Message).UpdatedAt,
		DeletedAt: &model.(*Message).DeletedAt,
	}
}
