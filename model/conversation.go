package model

import (
	"gochatbackend/pkg/databaseutil"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Conversation struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"  json:"id"`
	CreatedAt     time.Time          `bson:"createdAt"  json:"createdAt"`
	UpdatedAt     time.Time          `bson:"updatedAt"  json:"updatedAt"`
	DeletedAt     time.Time          `bson:"deletedAt"  json:"deletedAt"`
	Name          string             `bson:"name"  json:"name"`
	LastMessageAt time.Time          `bson:"lastMessageAt"  json:"lastMessageAt"`
}

func (u Conversation) GetBaseModel(model interface{}) databaseutil.BaseModel {
	return databaseutil.BaseModel{
		ID:        &model.(*Conversation).ID,
		CreatedAt: &model.(*Conversation).CreatedAt,
		UpdatedAt: &model.(*Conversation).UpdatedAt,
		DeletedAt: &model.(*Conversation).DeletedAt,
	}
}

type UserRole string

const (
	AdminRole  UserRole = "admin"
	MemberRole UserRole = "member"
)

type UserConversation struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CreatedAt      time.Time          `bson:"createdAt"  json:"createdAt"`
	UpdatedAt      time.Time          `bson:"updatedAt"  json:"updatedAt"`
	DeletedAt      time.Time          `bson:"deletedAt"  json:"deletedAt"`
	ConversationID primitive.ObjectID `bson:"conversationId" json:"conversationId"`
	UserID         primitive.ObjectID `bson:"userId" json:"userId"`
	Role           UserRole           `bson:"role" json:"role"`
}

func (u UserConversation) GetBaseModel(model interface{}) databaseutil.BaseModel {
	return databaseutil.BaseModel{
		ID:        &model.(*UserConversation).ID,
		CreatedAt: &model.(*UserConversation).CreatedAt,
		UpdatedAt: &model.(*UserConversation).UpdatedAt,
		DeletedAt: &model.(*UserConversation).DeletedAt,
	}
}
