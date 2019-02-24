package server

import (
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/mcavoyk/quirk/api/pkg/gfyid"

	"github.com/mcavoyk/quirk/api/pkg/ip"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	"github.com/mcavoyk/quirk/api/models"
)

type User struct {
	Username    string `json:"username" form:"username"`
	DisplayName string `json:"display_name" form:"display_name"`
	Password    string `json:"password" form:"password"`
	Email       string `json:"email" form:"email"`
}

const tokenExpiry = 14 * 24 * time.Hour

type Login struct {
	Username string  `json:"username" form:"username" binding:"required"`
	Password string  `json:"password" form:"password" binding:"required"`
	Lat      float64 `json:"lat" form:"lat" binding:"min=-90,max=90"`
	Lon      float64 `json:"lon" form:"lon" binding:"min=-180,max=180"`
}

func randomName() string {
	name := gfyid.RandomID()
	logrus.Infof("Random name: %s", name)
	return name
}

func marshalUser(user *User) *models.User {
	if user.DisplayName == "" {
		user.DisplayName = user.Username
	}
	newUser := &models.User{
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Email:       user.Email,
		Password:    hashSalt(user.Password),
	}
	newUser.ID = models.NewGUID()
	return newUser
}

func (env *Env) CreateUser(c *gin.Context) {
	newUser := new(User)
	_ = c.ShouldBind(newUser)

	logrus.Debugf("Got username: %s", newUser.Username)
	randomUser := false
	if newUser.Username == "" {
		randomUser = true
		newUser.Username = randomName()
	}

	existingUser := models.User{}
	err := env.db.ReadOne(&existingUser, models.SelectUserByName, newUser.Username)
	if err == nil {
		if randomUser {
			for err == nil {
				newUser.Username = randomName()
				err = env.db.ReadOne(&models.User{}, models.SelectUserByName, newUser.Username)
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"status": "Username already exists"})
			return
		}
	}

	randomPassword := false
	if newUser.Password == "" {
		randomPassword = true
		newUser.Password = NewPass()
	}

	err = env.db.Write(models.InsertUser, marshalUser(newUser))
	if err != nil {
		logrus.Errorf("Received Write Error: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
		return
	}

	user := models.User{}
	err = env.db.ReadOne(&user, models.SelectUserByName, newUser.Username)
	if err != nil {
		logrus.Errorf("Read user error: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": "internal server error"})
		return
	}
	user.Password = ""
	if randomPassword {
		user.Password = newUser.Password
	}
	c.JSON(http.StatusCreated, user)
}

func (env *Env) GetUser(c *gin.Context) {
	id := c.Param("id")
	user := models.User{}
	err := env.db.ReadOne(&user, models.SelectUser, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": NotFound})
		return
	}
	if err := env.HasPermission(c, c.GetString(UserContext), user.ID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"status": "Forbidden"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (env *Env) PatchUser(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString(UserContext)

	if err := env.HasPermission(c, userID, id); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"status": "Forbidden"})
		return
	}

	user := new(User)
	_ = c.ShouldBind(user)

	// Password change not supported at this time
	user.Password = ""

	if id == "" {
		id = userID
	}

	UserUpdate := &models.User{
		Default:     models.Default{ID: id},
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Password:    user.Password,
		Email:       user.Password,
	}
	err := env.db.Write(models.UpdateValues("users", *UserUpdate), UserUpdate.ID)

	if err != nil {
		logrus.Errorf("Error while updating user: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
		return
	}

	newUser := models.User{}
	err = env.db.ReadOne(&newUser, models.SelectUser)
	c.JSON(http.StatusOK, newUser)
}

func (env *Env) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString(UserContext)

	if err := env.HasPermission(c, userID, id); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"status": "Forbidden"})
		return
	}

	if _, err := env.db.Exec(models.DeleteUserSoft, id); err != nil {
		logrus.Errorf("Error while deleting user: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
		return
	}

	if _, err := env.db.Exec(models.DeleteSessions, id); err != nil {
		logrus.Errorf("Error while deleting user sessions: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (env *Env) LoginUser(c *gin.Context) {
	newLogin := new(Login)
	if err := c.ShouldBind(newLogin); err != nil {
		logrus.Debugf("Failed to bind to user struct: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Invalid or missing login fields",
		})
		return
	}

	user := models.User{}
	err := env.db.ReadOne(&user, models.SelectUserByName, newLogin.Username)
	if err != nil || user.DeletedAt != nil {
		logrus.Debugf("User %s not found on login", newLogin.Username)
		c.JSON(http.StatusUnauthorized, gin.H{"status": "Unauthorized"})
		return
	}
	validPass := comparePasswords(user.Password, newLogin.Password)

	if !validPass {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "Unauthorized"})
		return
	}

	session := &models.Session{
		UserID:    user.ID,
		Expiry:    time.Now().Add(tokenExpiry),
		Lat:       newLogin.Lat,
		Lon:       newLogin.Lon,
		IP:        ip.Parse(c.Request),
		UserAgent: c.GetHeader("User-Agent")}
	session.ID = models.NewGUID()

	err = env.db.Write(models.InsertSession, session)
	if err != nil {
		logrus.Errorf("Error inserting login sessions: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":  session.ID,
		"expiry": session.Expiry,
	})
}

// UserVerify is auth middleware to check session token from Authorization header
func (env *Env) UserVerify(c *gin.Context) {
	start := time.Now()
	sessionID := extractToken(c)
	if sessionID == "" {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	if val, ok := c.Get(RootKey); ok && val != nil {
		if val.(string) == sessionID {
			logrus.Warnf("RootKey used from ip [%s] for [%s: %s]", ip.Parse(c.Request), c.Request.Method, c.Request.URL.Path)
			c.Set(UserContext, sessionID)
			c.Next()
			return
		}
	}

	existingSession := models.Session{}
	err := env.db.ReadOne(&existingSession, models.SelectSession, sessionID)
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

	//FIXME: Might be too much to do a write on every request, possibly once per user per day?
	_ = env.db.Write(models.UpdateSession, existingSession)
	c.Set(UserContext, existingSession.UserID)
	logrus.Debugf("Total auth middleware: %f", time.Since(start).Seconds())
	c.Next()
}

func extractToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	headerSplit := strings.Split(authHeader, " ")
	if len(headerSplit) < 2 {
		return ""
	}
	return headerSplit[1]
}

func hashSalt(pwd string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		logrus.Error("Failed to hash password")
		panic("Failed to hash password")
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
