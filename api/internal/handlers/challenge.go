package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/ctf/api/internal/usecase"
	"github.com/gin-gonic/gin"
)

type ChallengeHandler struct {
	challengeUseCase usecase.ChallengeUsecase
}

func NewChallengeHandler(uc usecase.ChallengeUsecase) *ChallengeHandler {
	return &ChallengeHandler{
		challengeUseCase: uc,
	}
}

// CreateChallengeRequest defines the structure for challenge creation
var req struct {
	Name         string `json:"name" binding:"required"`
	Description  string `json:"description" binding:"required"`
	DownloadLink string `json:"downloadLink" binding:"required"`
	Hint         string `json:"hint" binding:"required"`
	Flag         string `json:"flag" binding:"required"`
	Score        string `json:"score" binding:"required"` // Keep score as a string first
}

// CreateChallenge handles the HTTP request to create a new challenge
func (h *ChallengeHandler) CreateChallenge(c *gin.Context) {

	// Bind JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert score from string to int
	score, err := strconv.Atoi(req.Score)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid score format"})
		return
	}
	// Call UseCase to create challenge
	err = h.challengeUseCase.CreateChallenge(req.Name, req.Description, req.DownloadLink, req.Hint, req.Flag, score)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Challenge created successfully"})
}

// SubmitFlagRequest defines the structure for flag submission
type SubmitFlagRequest struct {
	ChallengeID   int    `json:"challenge_id" binding:"required"`
	SubmittedFlag string `json:"submitted_flag" binding:"required"`
}

// SubmitFlag handles the HTTP request to submit a flag
func (h *ChallengeHandler) SubmitFlag(c *gin.Context) {
	var req SubmitFlagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	teamID, exists := c.Get("team_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: team_id not found in context"})
		return
	}

	teamIDInt, ok := teamID.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid team_id type in context"})
		return
	}

	success, err := h.challengeUseCase.VerifyAndSubmitFlag(teamIDInt, req.ChallengeID, req.SubmittedFlag)
	if err != nil {
		log.Printf("Error verifying flag: %v", err)

		if err.Error() == "team has already solved this challenge" {
			c.JSON(http.StatusConflict, gin.H{"message": "Flag already submitted"}) // HTTP 409 Conflict
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "An unexpected error occurred"})
		return
	}

	if !success {
		c.JSON(http.StatusForbidden, gin.H{"message": "Invalid flag", "success": false}) // HTTP 403 Forbidden
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Flag submitted successfully", "success": true}) // HTTP 200 OK
}

// GetChallenges retrieves and returns all challenges
func (h *ChallengeHandler) GetChallenges(c *gin.Context) {
	challenges, err := h.challengeUseCase.GetChallenges()
	log.Println(challenges)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, challenges)
}

// DeleteChallenge removes a challenge by name
func (h *ChallengeHandler) DeleteChallenge(c *gin.Context) {
	idParam := c.Param("id")
	log.Println("Received ID from URL:", idParam) // Debugging log

	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Println("Error converting ID:", err) // Debugging log
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid challenge ID"})
		return
	}

	log.Println("Parsed ID:", id)

	err = h.challengeUseCase.DeleteChallenge(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Challenge deleted successfully"})
}
