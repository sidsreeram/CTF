package handlers

import (
	"net/http"

	"github.com/ctf/api/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Handlers struct {
	uc usecase.UseCase
}

func NewHandlers(uc usecase.UseCase) *Handlers {
	return &Handlers{uc: uc}
}

func (h *Handlers) GetTeams(c *gin.Context) {
	teams, err := h.uc.GetTeams()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	c.JSON(http.StatusOK, teams)
}

func (h *Handlers) GetChallenges(c *gin.Context) {
	challenges, err := h.uc.GetChallenges() // Removed context parameter
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	c.JSON(http.StatusOK, challenges)
}

func (h *Handlers) GetScores(c *gin.Context) {
    teamscore, err := h.uc.GetScores() // Fetch scores
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
        return
    }

    // Ensure the response is an array of objects with proper field names
    c.JSON(http.StatusOK, gin.H{"scores": teamscore})
}



var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (h *Handlers) HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	for {
		var data struct {
			TeamName      string `json:"team_name"`
			ChallengeName string `json:"challenge_name"`
			Score         int    `json:"score"`
		}

		if err := conn.ReadJSON(&data); err != nil {
			break
		}

		if data.TeamName == "" || data.ChallengeName == "" || data.Score < 0 {
			continue
		}

		if err := h.uc.UpdateScore(data.TeamName, data.ChallengeName, data.Score); err != nil { // Removed context
			continue
		}
	}
}
