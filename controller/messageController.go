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

type MessageController interface {
	PostMessages(ctx *gin.Context)
	GetMessages(ctx *gin.Context)
}

type messageController struct {
	messageService service.MessageService
}

func NewMessageController(messageService service.MessageService) MessageController {
	return messageController{
		messageService: messageService,
	}
}

func (m messageController) PostMessages(ctx *gin.Context) {
	var request model.SendMessageForm
	if err := ctx.ShouldBindJSON(&request); err != nil {
		log.Errorf("invalid post message request: %s", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, model.Response("invalid request"))
		return
	}
	userId, err := auth.ParseIdFromCtx(ctx)
	if err != nil {
		log.Errorf("failed to parse user id from context: %s", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, model.Response("invalid jwt"))
		return
	}
	if request.Content == "" {
		log.Errorf("invalid request content: %s", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, model.Response("invalid request"))
		return
	}
	objectConversationId, err := primitive.ObjectIDFromHex(request.ConversationId)
	if err != nil {
		log.Errorf("failed to parse conversation id to object id: %s", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, model.Response("invalid request"))
		return
	}
	message := model.Message{
		Content:        request.Content,
		UserID:         userId,
		ConversationID: objectConversationId,
	}
	result, err := m.messageService.CreateMessage(ctx, message)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, model.Response("failed to send new message"))
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (m messageController) GetMessages(ctx *gin.Context) {
	limitStr := ctx.Query("limit")
	limit, _ := strconv.ParseInt(limitStr, 10, 64)
	conversationId := ctx.Query("conversationId")
	objectConversationId, err := primitive.ObjectIDFromHex(conversationId)
	if err != nil {
		log.Errorf("failed to parse conversation id to object id: %s", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, model.Response("invalid request"))
		return
	}
	result, err := m.messageService.GetMessageConversationId(ctx, objectConversationId, limit)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, model.Response("failed to get message"))
		return
	}
	ctx.JSON(http.StatusOK, result)
}
