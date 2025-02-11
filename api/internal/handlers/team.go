package handlers

import (
	"net/http"
	"strconv"

	"github.com/ctf/api/internal/usecase"
	"github.com/gin-gonic/gin"
)

type TeamHandler struct {
	usecase usecase.TeamUsecase
}

func NewTeamHandler(uc usecase.TeamUsecase) *TeamHandler {
	return &TeamHandler{usecase: uc}

}

func (h *TeamHandler) RegisterTeam(c *gin.Context) {
	var req struct {
		TeamName string `json:"team_name"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	err := h.usecase.RegisterTeam(req.TeamName, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Team registered successfully"})
}

// var jwtSecret = []byte("secret")

func (h *TeamHandler) LoginTeam(c *gin.Context) {
	var req struct {
		TeamName string `json:"team_name"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	token, role, teamID, err := h.usecase.LoginTeam(req.TeamName, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Set authentication cookies
	if role == "admin" {
		c.SetCookie("AdminAuth", token, 3600, "/", "", false, true)
	} else {
		c.SetCookie("AuthToken", token, 3600, "/", "", false, true)
	}

	c.SetCookie("TeamID", strconv.Itoa(teamID), 3600, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message":       "Login successful",
		"authenticated": true,
		"role":          role,
		"team_id":       teamID,
	})
}
func (h *TeamHandler) BlockTeam(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}

	err = h.usecase.BlockTeam(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to block team"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Team blocked successfully"})
}
func (h *TeamHandler) UnblockTeam(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}

	err = h.usecase.UnblockTeam(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unblock team"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Team unblocked successfully"})
}
