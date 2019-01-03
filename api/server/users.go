package server

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mcavoyk/quirk/ip"
	"github.com/mcavoyk/quirk/models"
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
	existingUser.UsedAt = time.Now()
	env.DB.UserUpdate(existingUser)
	c.Next()
}

func (env *Env) UserCreate(c *gin.Context) {
	userID := env.DB.UserInsert(&models.User{IP: ip.Parse(c.Request)})
	c.JSON(http.StatusOK, gin.H{
		"token": userID,
	})
}

func (env *Env) UserValidate(c *gin.Context) {
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
