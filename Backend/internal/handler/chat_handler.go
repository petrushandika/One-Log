package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/petrushandika/one-log/internal/service"
	"github.com/petrushandika/one-log/pkg/utils"
)

type ChatHandler struct {
	service service.ChatService
}

func NewChatHandler(svc service.ChatService) *ChatHandler {
	return &ChatHandler{service: svc}
}

type chatRequest struct {
	Message string `json:"message" binding:"required"`
}

// Ask handles POST /api/chat
func (h *ChatHandler) Ask(c *gin.Context) {
	var req chatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "Message is required", err.Error())
		return
	}

	userID := c.GetUint("user_id")
	reply, err := h.service.Ask(req.Message, userID)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "AI service unavailable", err.Error())
		return
	}

	utils.Success(c, http.StatusOK, "Reply generated", gin.H{"reply": reply})
}
