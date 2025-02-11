package handlers

import (
	"net/http"

	"github.com/ctf/api/internal/models"
	"github.com/gin-gonic/gin"
)

type WebSocketHandler struct {
	webSocketService models.WebSocketService
}

func NewWebSocketHandler(ws models.WebSocketService) *WebSocketHandler {
	return &WebSocketHandler{webSocketService: ws}
}

// HTTP handler to trigger score updates
func (h *WebSocketHandler) UpdateScore(c *gin.Context) {
	var request struct {
		ChallengeID int `json:"challenge_id"`
		NewScore    int `json:"new_score"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Broadcast the score update to WebSocket clients
	h.webSocketService.BroadcastScoreUpdate(request.ChallengeID, request.NewScore)

	c.JSON(http.StatusOK, gin.H{"message": "Score updated successfully"})
}
