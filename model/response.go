package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type implResponse struct {
	Message string `json:"message"`
}

func Response(message string) implResponse {
	return implResponse{
		Message: message,
	}
}

type MessageResponse struct {
	ID             primitive.ObjectID `bson:"_id" json:"id"`
	CreatedAt      time.Time          `bson:"createdAt"  json:"createdAt"`
	UpdatedAt      time.Time          `bson:"updatedAt"  json:"updatedAt"`
	DeletedAt      time.Time          `bson:"deletedAt"  json:"deletedAt"`
	Content        string             `bson:"content"  json:"content"`
	UserID         primitive.ObjectID `bson:"userId"  json:"userId"`
	ConversationID primitive.ObjectID `bson:"conversationId" json:"conversationId"`
	UserDetail     User               `bson:"userDetail"  json:"userDetail"`
}

type ConversationResponse struct {
	ID            primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	CreatedAt     time.Time            `bson:"createdAt"  json:"createdAt"`
	UpdatedAt     time.Time            `bson:"updatedAt"  json:"updatedAt"`
	DeletedAt     time.Time            `bson:"deletedAt"  json:"deletedAt"`
	Name          string               `bson:"name"  json:"name"`
	LastMessageAt time.Time            `bson:"lastMessageAt"  json:"lastMessageAt"`
	Members       []primitive.ObjectID `bson:"members" json:"members"`
}
