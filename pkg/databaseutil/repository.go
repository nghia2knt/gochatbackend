package databaseutil

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BaseModel struct {
	ID        *primitive.ObjectID
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}

type BaseModelI interface {
	GetBaseModel(pointer interface{}) BaseModel
}

type BaseRepository[M BaseModelI] interface {
	Create(ctx context.Context, db *mongo.Database, model *M) (*M, error)
	Update(ctx context.Context, db *mongo.Database, id primitive.ObjectID, model map[string]interface{}) error
	Delete(ctx context.Context, db *mongo.Database, id primitive.ObjectID) error
	Find(ctx context.Context, db *mongo.Database, filter map[string]interface{}, options ...*options.FindOptions) ([]*M, error)
	FindOne(ctx context.Context, db *mongo.Database, filter map[string]interface{}) (*M, error)
	FindByID(ctx context.Context, db *mongo.Database, id primitive.ObjectID) (*M, error)
	Count(ctx context.Context, db *mongo.Database, filter map[string]interface{}) (int64, error)
}

type baseRepository[M BaseModelI] struct {
	collection string
}

func NewBaseRepository[M BaseModelI](collection string) BaseRepository[M] {
	return baseRepository[M]{
		collection: collection,
	}
}

func (b baseRepository[M]) Create(ctx context.Context, db *mongo.Database, model *M) (*M, error) {
	baseModel := (*model).GetBaseModel(model)
	createdAt := time.Now()
	*baseModel.CreatedAt = createdAt
	result, err := db.Collection(b.collection).InsertOne(ctx, model)
	if err != nil {
		return nil, err
	}
	*baseModel.ID = result.InsertedID.(primitive.ObjectID)
	return model, nil
}

func (b baseRepository[M]) Update(ctx context.Context, db *mongo.Database, id primitive.ObjectID, model map[string]interface{}) error {
	var findModel M
	err := db.Collection(b.collection).FindOne(ctx, bson.M{"_id": id}).Decode(&findModel)
	if err != nil {
		return err
	}
	deletedAt := findModel.GetBaseModel(&findModel).DeletedAt
	if !deletedAt.IsZero() {
		return fmt.Errorf("not found model")
	}
	model["updatedAt"] = time.Now()
	_, err = db.Collection(b.collection).UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": model})
	if err != nil {
		return err
	}
	return nil
}

func (b baseRepository[M]) Delete(ctx context.Context, db *mongo.Database, id primitive.ObjectID) error {
	var findModel M
	err := db.Collection(b.collection).FindOne(ctx, bson.M{"_id": id}).Decode(&findModel)
	if err != nil {
		return err
	}
	deletedAt := findModel.GetBaseModel(&findModel).DeletedAt
	if !deletedAt.IsZero() {
		return fmt.Errorf("not found model")
	}
	_, err = db.Collection(b.collection).UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": map[string]interface{}{
		"deletedAt": time.Now(),
	}})
	if err != nil {
		return err
	}
	return nil
}

func (b baseRepository[M]) Find(ctx context.Context, db *mongo.Database, filter map[string]interface{}, options ...*options.FindOptions) ([]*M, error) {
	var results []*M
	var deletedAt time.Time
	filter["deletedAt"] = deletedAt
	cursor, err := db.Collection(b.collection).Find(ctx, filter, options...)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var result M
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		results = append(results, &result)
	}
	return results, nil
}

func (b baseRepository[M]) FindOne(ctx context.Context, db *mongo.Database, filter map[string]interface{}) (*M, error) {
	var result M
	var deletedAt time.Time
	filter["deletedAt"] = deletedAt
	err := db.Collection(b.collection).FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (b baseRepository[M]) FindByID(ctx context.Context, db *mongo.Database, id primitive.ObjectID) (*M, error) {
	var result M
	var deletedAt time.Time
	err := db.Collection(b.collection).FindOne(ctx, bson.M{"_id": id, "deletedAt": deletedAt}).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (b baseRepository[M]) Count(ctx context.Context, db *mongo.Database, filter map[string]interface{}) (int64, error) {
	var deletedAt time.Time
	filter["deletedAt"] = deletedAt
	count, err := db.Collection(b.collection).CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}
	return count, nil
}
