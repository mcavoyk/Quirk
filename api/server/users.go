package server

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mcavoyk/quirk/api/gfyid"
	"github.com/mcavoyk/quirk/api/location"
	"github.com/mcavoyk/quirk/api/models"
)

type User struct {
	Name     string  `json:"name" form:"name"`
	Password string  `json:"-" form:"password"`
	Email    string  `json:"email" form:"email"`
	Lat      float64 `json:"lat" form:"lat" binding:"required"`
	Lon      float64 `json:"lon" form:"lon" binding:"required"`
}

func (env *Env) UserVerify(c *gin.Context) {
	userID := extractToken(c)
	if userID == "" {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	existingUser := env.DB.GetUser(userID)
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

	env.DB.UserUpdate(existingUser)
	c.Set(UserContext, existingUser.ID)
	c.Next()
}

func (env *Env) CreateUser(c *gin.Context) {
	newUser := new(User)
	if err := c.ShouldBind(newUser); err != nil {
		env.Log.Debugf("Failed to bind to user struct: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Missing or invalid latitude and longitude",
		})
		return
	}

	if newUser.Name == "" {
		newUser.Name = gfyid.RandomID()
	}

	user := env.DB.InsertUser(&models.User{
		Name: newUser.Name,
		IP:   c.ClientIP(),
		Lat:  location.ToRadians(newUser.Lat),
		Lon:  location.ToRadians(newUser.Lon)})
	c.JSON(http.StatusCreated, user)
}

func (env *Env) ValidateUser(c *gin.Context) {
	userID := c.Param("token")
	if userID == "" {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	existingUser := env.DB.GetUser(userID)
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
