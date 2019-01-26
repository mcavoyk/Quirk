package server

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mcavoyk/quirk/api/ip"
	"github.com/mcavoyk/quirk/api/models"
)

func (env *Env) UserVerify(c *gin.Context) {
	userID := extractToken(c)
	if userID == "" {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	existingUser := env.DB.UserGet(userID)
	if existingUser.ID != userID {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	// Extract and use coords if they are included and valid
	coords, err := extractCoords(c)
	if err == nil {
		existingUser.Lat = coords.Lat
		existingUser.Lon = coords.Lon
	}

	existingUser.UsedAt = time.Now()
	env.DB.UserUpdate(existingUser)
	c.Set(UserContext, existingUser.ID)
	c.Next()
}

func (env *Env) CreateUser(c *gin.Context) {
	coords, err := extractCoords(c)

	if err != nil  {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	userID := env.DB.UserInsert(&models.User{IP: ip.Parse(c.Request), Lat: coords.Lat, Lon: coords.Lon})
	c.JSON(http.StatusOK, gin.H{
		"token": userID,
	})
}

func (env *Env) ValidateUser(c *gin.Context) {
	userID := c.Param("token")
	if userID == "" {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	existingUser := env.DB.UserGet(userID)
	if existingUser.ID != userID {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	c.JSON(http.StatusOK, existingUser)
}

func extractToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	headerSplit := strings.Split(authHeader, " ")
	if len(headerSplit) < 2 {
		return ""
	}
	return headerSplit[1]
}
