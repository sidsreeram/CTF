package handlers

import (
	"log"
	"net/http"
	"text/template"

	"github.com/ctf-api/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Handlers struct {
	uc *usecase.Usecase
}

func NewHandlers(uc *usecase.Usecase) *Handlers {
	return &Handlers{uc: uc}
}

func (h *Handlers) GetTeams(c *gin.Context) {
	teams, err := h.uc.GetTeams() 
	log.Println("hanlder reached")
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
	teamscore, err := h.uc.GetScores() // Removed context parameter
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	funcMap := template.FuncMap{
		"add": func(a, b int) int { return a + b },
	}
	tmpl, err := template.New("hackerboard.html").Funcs(funcMap).ParseFiles("../../../template/hackerboard.html")
	if err != nil {
		log.Println("Error loading template:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error loading template"})
		return
	}
	err = tmpl.Execute(c.Writer, gin.H{"Scores": teamscore})
	if err != nil {
		log.Println("Error executing template:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error executing template"})
		return
	}
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
