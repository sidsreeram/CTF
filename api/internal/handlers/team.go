package handlers

import (
	"net/http"

	"github.com/ctf-api/internal/usecase"
	"github.com/gin-gonic/gin"
)

type TeamHandler struct {
	usecase usecase.TeamUsecase
}

func NewTeamHandler( uc usecase.TeamUsecase)*TeamHandler{
return  &TeamHandler{usecase: uc}

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

func (h *TeamHandler) LoginTeam(c *gin.Context) {
	var req struct {
		TeamName string `json:"team_name"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	success, err := h.usecase.LoginTeam(req.TeamName, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "authenticated": success})
}
