package middleware

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func errorHandler(c *gin.Context, err error) {
	if _, ok := err.(*InvalidTokenError); ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid token",
		})
		return
	}

	// Handle other errors here
}
func UserAuth(c *gin.Context) {
	// s := c.Request.Header.Get("Authorization")
	tokenString, err := c.Cookie("UserAuth")
	if err != nil {
		errorHandler(c, err)
		return
	}

	userId, err := ValidateToken(tokenString)
	if err != nil {
		errorHandler(c, err)
		return
	}

	userIdStr := strconv.Itoa(userId)
	c.SetCookie("userId", userIdStr, 3600, "/", "", false, true)

	c.Next()
}
