package controller

import (
	"gochatbackend/model"
	"gochatbackend/pkg/auth"
	"gochatbackend/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type UserController interface {
	GetUsers(c *gin.Context)
	PostUsers(c *gin.Context)
	Login(c *gin.Context)
	GetIdentity(c *gin.Context)
}

type userController struct {
	userServcie service.UserService
}

func NewUserController(userService service.UserService) UserController {
	return userController{
		userServcie: userService,
	}
}

func (u userController) GetUsers(c *gin.Context) {
	limitStr := c.Query("limit")
	limit, _ := strconv.ParseInt(limitStr, 10, 64)
	skipStr := c.Query("skip")
	skip, _ := strconv.ParseInt(skipStr, 10, 64)
	username := c.Query("username")
	result, err := u.userServcie.GetUser(c, username, limit, skip)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.Response("failed to get user"))
		return
	}
	c.JSON(http.StatusOK, result)
}

func (u userController) PostUsers(c *gin.Context) {
	var request model.CreateUserForm
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Errorf("invalid create user request: %s", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, model.Response("invalid request"))
		return
	}
	user := model.User{
		Name:     request.Name,
		Username: request.Username,
		Password: request.Password,
	}
	_, err := u.userServcie.CreateUser(c, user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.Response("failed to create user"))
		return
	}
	c.JSON(http.StatusOK, model.Response("ok"))
}

func (u userController) Login(c *gin.Context) {
	var request model.LoginForm
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Errorf("invalid login request: %s", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, model.Response("invalid request"))
		return
	}
	token, err := u.userServcie.Login(c, request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.Response("failed to login"))
		return
	}
	c.JSON(http.StatusOK, model.Response(token))
}

func (u userController) GetIdentity(c *gin.Context) {
	userId, err := auth.ParseIdFromCtx(c)
	if err != nil {
		log.Errorf("failed to parrse user id from context: %s", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, model.Response("invalid jwt"))
		return
	}
	user, err := u.userServcie.GetUserById(c, userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.Response("failed to get user"))
		return
	}
	c.JSON(http.StatusOK, user)
}
