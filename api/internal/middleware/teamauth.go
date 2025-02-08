package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func TeamAuth(c *gin.Context) {
	tokenString, err := c.Cookie("AuthToken")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authentication token"})
		c.Abort()
		return
	}

	teamID, err := ValidateToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authentication token"})
		c.Abort()
		return
	}

	// Store team ID in context
	c.Set("team_id", teamID)
	c.Next()
}