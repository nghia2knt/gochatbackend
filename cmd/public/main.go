package main

import (
	"context"
	"gochatbackend/controller"
	"gochatbackend/repository"
	"gochatbackend/service"

	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DATABASE_URI  = "mongodb://host.docker.internal:27017"
	DATABASE_NAME = "test"
)

func main() {
	r := gin.Default()
	db, _ := database()
	ctx := context.TODO()
	userController := initUserController(ctx, db)
	messageController := initMessageController(ctx, db)
	conversation := initConversationController(ctx, db)
	r.Use(CORSMiddleware())

	r.POST("/register", userController.PostUsers)
	r.POST("/login", userController.Login)

	r.GET("/users", userController.GetUsers)
	r.GET("/identity", userController.GetIdentity)

	r.GET("/messages", messageController.GetMessages)
	r.POST("/messages", messageController.PostMessages)

	r.GET("/conversations", conversation.GetIdentityConversation)
	r.POST("/conversations", conversation.PostConversation)
	r.GET("/conversations/:conversationId", conversation.GetConversationById)

	r.Run(":9010")
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, POST")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func database() (*mongo.Database, *mongo.Client) {
	clientOptions := options.Client().ApplyURI(DATABASE_URI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Hello MDFK, Connected to MongoDB!")
	return client.Database(DATABASE_NAME), client
}

func initUserController(ctx context.Context, db *mongo.Database) controller.UserController {
	userRepository := repository.NewUserRepository()
	userService := service.NewUserService(db, userRepository)
	err := userRepository.CreateIndex(ctx, db)
	if err != nil {
		fmt.Println(err)
	}
	return controller.NewUserController(userService)
}

func initMessageController(ctx context.Context, db *mongo.Database) controller.MessageController {
	messageRepository := repository.NewMessageRepository()
	conversationRepository := repository.NewConversationRepository()
	err := messageRepository.CreateIndex(ctx, db)
	if err != nil {
		fmt.Println(err)
	}
	userRepository := repository.NewUserRepository()
	messageService := service.NewMessageService(db, messageRepository, userRepository, conversationRepository)
	return controller.NewMessageController(messageService)
}

func initConversationController(ctx context.Context, db *mongo.Database) controller.ConversationController {
	conversationRepository := repository.NewConversationRepository()
	err := conversationRepository.CreateIndex(ctx, db)
	if err != nil {
		fmt.Println(err)
	}
	conversationService := service.NewConversationService(db, conversationRepository)
	return controller.NewConversationController(conversationService)
}
