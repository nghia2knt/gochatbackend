package service

import (
	"context"
	"fmt"
	"gochatbackend/model"
	"gochatbackend/pkg/auth"
	"gochatbackend/repository"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserService interface {
	CreateUser(ctx context.Context, user model.User) (model.User, error)
	GetUser(ctx context.Context, username string, limit int64, skip int64) ([]model.User, error)
	Login(ctx context.Context, loginRequest model.LoginForm) (string, error)
	GetUserById(ctx context.Context, id primitive.ObjectID) (model.User, error)
}

type userService struct {
	db             *mongo.Database
	userRepository repository.UserRepository
}

func NewUserService(
	db *mongo.Database,
	userRepository repository.UserRepository,
) UserService {
	return userService{
		db:             db,
		userRepository: userRepository,
	}
}

func (u userService) CreateUser(ctx context.Context, user model.User) (model.User, error) {
	bytePassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		log.Error(err)
		return model.User{}, err
	}
	user.Password = string(bytePassword)
	result, err := u.userRepository.Create(ctx, u.db, &user)
	if err != nil {
		log.Error(err)
		return model.User{}, err
	}
	return *result, nil
}

func (u userService) GetUser(ctx context.Context, username string, limit int64, skip int64) ([]model.User, error) {
	if limit == 0 {
		limit = 20
	}
	filter := bson.M{}
	if username != "" {
		filter["username"] = bson.M{
			"$regex": ".*(?i)" + username + "(?i).*",
		}
	}
	result, err := u.userRepository.Find(ctx, u.db, filter, options.Find().SetLimit(limit).SetSort(map[string]interface{}{"createdAt": -1}).SetSkip(skip))
	if err != nil {
		log.Error(err)
		return nil, err
	}
	var userList []model.User
	for _, value := range result {
		userList = append(userList, *value)
	}
	return userList, nil
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (u userService) Login(ctx context.Context, loginRequest model.LoginForm) (string, error) {
	user, err := u.userRepository.FindOne(ctx, u.db, map[string]interface{}{
		"username": loginRequest.Username,
	})
	if err != nil {
		log.Error(err)
		return "", err
	}
	check := checkPasswordHash(loginRequest.Password, user.Password)
	if !check {
		return "", fmt.Errorf("check password error")
	}
	token, err := auth.GenerateToken(*user)
	if err != nil {
		log.Error(err)
		return "", err
	}
	return token, nil
}

func (u userService) GetUserById(ctx context.Context, id primitive.ObjectID) (model.User, error) {
	result, err := u.userRepository.FindByID(ctx, u.db, id)
	if err != nil {
		log.Error(err)
		return model.User{}, err
	}
	return *result, nil
}
