package repository

import (
	"context"
	"gochatbackend/model"
	"gochatbackend/pkg/databaseutil"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository interface {
	databaseutil.BaseRepository[model.User]
	CreateIndex(ctx context.Context, db *mongo.Database) error
}

type userRepository struct {
	databaseutil.BaseRepository[model.User]
}

func NewUserRepository() UserRepository {
	return userRepository{databaseutil.NewBaseRepository[model.User]("users")}
}

func (u userRepository) CreateIndex(ctx context.Context, db *mongo.Database) error {
	_, err := db.Collection("users").Indexes().CreateOne(
		ctx,
		mongo.IndexModel{
			Keys:    bson.M{"username": 1},
			Options: options.Index().SetName("username_index").SetUnique(true),
		},
	)
	if err != nil {
		return err
	}
	return nil
}
