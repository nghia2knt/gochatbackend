package main

import (
	"context"
	"gochatbackend/controller"
	"gochatbackend/pkg"
	"gochatbackend/repository"
	"gochatbackend/service"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	r := gin.Default()
	db, _ := pkg.Database()
	ctx := context.TODO()
	messageController := initMessageController(ctx, db)
	r.Use(pkg.CORSMiddleware())
	r.GET("/messages", messageController.GetMessages)
	r.POST("/messages", messageController.PostMessages)
	r.Run(":9003")
}

func initMessageController(ctx context.Context, db *mongo.Database) controller.MessageController {
	messageRepository := repository.NewMessageRepository()
	conversationRepository := repository.NewConversationRepository()
	err := messageRepository.CreateIndex(ctx, db)
	if err != nil {
		log.Fatal(err)
	}
	userRepository := repository.NewUserRepository()
	messageService := service.NewMessageService(db, messageRepository, userRepository, conversationRepository)
	return controller.NewMessageController(messageService)
}
