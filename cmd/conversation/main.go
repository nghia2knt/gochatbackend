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
	conversation := initConversationController(ctx, db)
	r.Use(pkg.CORSMiddleware())
	r.GET("/conversations", conversation.GetIdentityConversation)
	r.POST("/conversations", conversation.PostConversation)
	r.GET("/conversations/:conversationId", conversation.GetConversationById)
	r.Run(":9002")
}

func initConversationController(ctx context.Context, db *mongo.Database) controller.ConversationController {
	conversationRepository := repository.NewConversationRepository()
	err := conversationRepository.CreateIndex(ctx, db)
	if err != nil {
		log.Fatal(err)
	}
	conversationService := service.NewConversationService(db, conversationRepository)
	return controller.NewConversationController(conversationService)
}
