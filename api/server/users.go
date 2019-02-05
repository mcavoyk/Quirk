package server

import (
	"net/http"
	"strings"
	"time"

	"github.com/mcavoyk/quirk/api/pkg/ip"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	"github.com/mcavoyk/quirk/api/models"
	"github.com/mcavoyk/quirk/api/pkg/gfyid"
)

type User struct {
	Username    string `json:"username" form:"username"`
	DisplayName string `json:"display_name" form:"display_name"`
	Password    string `json:"password" form:"password"`
	Email       string `json:"email" form:"email"`
}

type Login struct {
	Username string  `json:"username" form:"username" binding:"required"`
	Password string  `json:"password" form:"password" binding:"required"`
	Lat      float64 `json:"lat" form:"lat" binding:"required"`
	Lon      float64 `json:"lon" form:"lon" binding:"required"`
}

func (env *Env) UserVerify(c *gin.Context) {
	sessionID := extractToken(c)
	if sessionID == "" {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	existingSession, err := env.DB.GetSession(sessionID)
	if err != nil {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	existingUser, err := env.DB.GetUserBySession(sessionID)
	if err != nil {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	// Extract and use coords if they are included and valid
	coords, err := extractCoords(c)
	if err == nil {
		existingSession.Lat = coords.Lat
		existingSession.Lon = coords.Lon
	}
	existingSession.IP = ip.Parse(c.Request)

	env.DB.SessionUpdate(existingSession)
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

	if newUser.Username != "" {
		_, err := env.DB.GetUserByName(newUser.Username)
		if err == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "Username already exists",
			})
			return
		}
	} else {
		newUser.Username = gfyid.RandomID()
		_, err := env.DB.GetUserByName(newUser.Username)
		for err != nil {
			newUser.Username = gfyid.RandomID()
			_, err = env.DB.GetUserByName(newUser.Username)
		}
	}

	if newUser.DisplayName == "" {
		newUser.DisplayName = newUser.Username
	}

	randomPassword := false
	if newUser.Password == "" {
		randomPassword = true
		newUser.Password = NewPass()
	}

	user, err := env.DB.InsertUser(&models.User{
		DisplayName: newUser.DisplayName,
		Password:    env.hashAndSalt(newUser.Password),
		Email:       newUser.Email,
		Username:    newUser.Username,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": err.Error(),
		})
		return
	}

	if randomPassword {
		user.Password = newUser.Password
	} else {
		user.Password = ""
	}

	c.JSON(http.StatusCreated, user)
}

func (env *Env) LoginUser(c *gin.Context) {
	newLogin := new(Login)
	if err := c.ShouldBind(newLogin); err != nil {
		env.Log.Debugf("Failed to bind to user struct: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Invalid or missing login fields",
		})
		return
	}

	user, err := env.DB.GetUserByName(newLogin.Username)
	if err != nil {
		env.Log.Debugf("User %s not found on login", newLogin.Username)
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "Unauthorized",
		})
		return
	}

	validPass := comparePasswords(user.Password, newLogin.Password)
	if validPass {
		expiry := time.Now().Add(14 * 24 * time.Hour)
		session, err := env.DB.InsertSession(&models.Session{
			UserID:    user.ID,
			Expiry:    expiry,
			Lat:       newLogin.Lat,
			Lon:       newLogin.Lon,
			IP:        ip.Parse(c.Request),
			UserAgent: c.GetHeader("User-Agent"),
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token":   session.ID,
			"expires": session.Expiry,
		})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "Unauthorized",
		})
	}
}

func (env *Env) GetUser(c *gin.Context) {
	id := c.Param("id")
	user, err := env.DB.GetUser(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "Page not found",
		})
		return
	}

	if user.ID != c.GetString(UserContext) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "Unauthorized",
		})
		return
	}

	user.Password = ""
	c.JSON(http.StatusOK, user)
}

func (env *Env) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString(UserContext)

	if err := env.HasPermission(userID, id); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"status": "Forbidden"})
		return
	}

	_ = env.DB.DeleteUser(id)
	c.Status(http.StatusNoContent)
}

func extractToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	headerSplit := strings.Split(authHeader, " ")
	if len(headerSplit) < 2 {
		return ""
	}
	return headerSplit[1]
}

func (env *Env) hashAndSalt(pwd string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		env.Log.Error("Failed to hash password")
		return pwd
	}
	return string(hash)
}

func comparePasswords(hashedPwd, plainPwd string) bool {
	byteHash := []byte(hashedPwd)
	bytePlain := []byte(plainPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePlain)
	if err != nil {
		return false
	}

	return true
}
