package controller

import (
	"gochatbackend/model"
	"gochatbackend/pkg/auth"
	"gochatbackend/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ConversationController interface {
	PostConversation(ctx *gin.Context)
	GetIdentityConversation(ctx *gin.Context)
	GetConversationById(ctx *gin.Context)
}

type conversationController struct {
	conversationService service.ConversationService
}

func NewConversationController(conversationService service.ConversationService) ConversationController {
	return conversationController{
		conversationService: conversationService,
	}
}

func (c conversationController) PostConversation(ctx *gin.Context) {
	var request model.CreateConversationForm
	if err := ctx.ShouldBindJSON(&request); err != nil {
		log.Errorf("invalid post conversation request: %s", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, model.Response("invalid request"))
		return
	}
	userId, err := auth.ParseIdFromCtx(ctx)
	if err != nil {
		log.Errorf("failed to parse userid from context: %s", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, model.Response("invalid jwt"))
		return
	}
	var members []primitive.ObjectID
	for _, value := range request.Members {
		memberId, err := primitive.ObjectIDFromHex(value)
		if err != nil {
			log.Errorf("failed to parse member id to object id: %s", err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, model.Response("invalid request"))
			return
		}
		members = append(members, memberId)
	}
	result, err := c.conversationService.CreateConversation(ctx, request.Name, userId, members)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, model.Response("failed to create new conversation"))
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (c conversationController) GetIdentityConversation(ctx *gin.Context) {
	userId, err := auth.ParseIdFromCtx(ctx)
	if err != nil {
		log.Errorf("failed to parse userid from context: %s", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, model.Response("invalid jwt"))
		return
	}
	limitStr := ctx.Query("limit")
	limit, _ := strconv.ParseInt(limitStr, 10, 64)
	result, err := c.conversationService.GetConversationUserId(ctx, userId, limit)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, model.Response("failed to get conversations"))
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (c conversationController) GetConversationById(ctx *gin.Context) {
	userId, err := auth.ParseIdFromCtx(ctx)
	if err != nil {
		log.Errorf("failed to parse userid from context: %s", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, model.Response("invalid jwt"))
		return
	}
	conversationIdStr := ctx.Param("conversationId")
	conversationId, err := primitive.ObjectIDFromHex(conversationIdStr)
	if err != nil {
		log.Errorf("failed to parse conversation id to object id: %s", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, model.Response("invalid request"))
		return
	}
	result, err := c.conversationService.GetConversationById(ctx, userId, conversationId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, model.Response("failed to get conversation"))
		return
	}
	ctx.JSON(http.StatusOK, result)
}
