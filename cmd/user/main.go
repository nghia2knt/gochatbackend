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
	userController := initUserController(ctx, db)
	r.Use(pkg.CORSMiddleware())
	r.POST("/register", userController.PostUsers)
	r.POST("/login", userController.Login)
	r.GET("/users", userController.GetUsers)
	r.GET("/identity", userController.GetIdentity)
	r.Run(":9001")
}

func initUserController(ctx context.Context, db *mongo.Database) controller.UserController {
	userRepository := repository.NewUserRepository()
	userService := service.NewUserService(db, userRepository)
	err := userRepository.CreateIndex(ctx, db)
	if err != nil {
		log.Fatal(err)
	}
	return controller.NewUserController(userService)
}
