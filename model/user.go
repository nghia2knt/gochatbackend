package model

import (
	"gochatbackend/pkg/databaseutil"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time          `bson:"createdAt"  json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt"  json:"updatedAt"`
	DeletedAt time.Time          `bson:"deletedAt"  json:"deletedAt"`
	Name      string             `bson:"name"  json:"name"`
	Username  string             `bson:"username"  json:"username"`
	Age       int64              `bson:"age" json:"age"`
	Password  string             `bson:"password" json:"password"`
}

func (u User) GetBaseModel(model interface{}) databaseutil.BaseModel {
	return databaseutil.BaseModel{
		ID:        &model.(*User).ID,
		CreatedAt: &model.(*User).CreatedAt,
		UpdatedAt: &model.(*User).UpdatedAt,
		DeletedAt: &model.(*User).DeletedAt,
	}
}
