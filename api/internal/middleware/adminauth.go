package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminAuth(c *gin.Context) {
	tokenString, err := c.Cookie("AdminAuth")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "admin token required"})
		c.Abort()
		return
	}

	adminID,  err := ValidateToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authentication token"})
		c.Abort()
		return
	}


	c.Set("adminID", adminID)
	c.Next()
}
