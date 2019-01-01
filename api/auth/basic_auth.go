package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ExtractUser(c *gin.Context) string {
	return c.GetHeader("User")
}

func VerifyUser(c *gin.Context) {
	user := ExtractUser(c)
	// TODO: Verify user ID is in the proper GUID format
	if user == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid User authentication",
		})
		return
	}
	c.Next()
}