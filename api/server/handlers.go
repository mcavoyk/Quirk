package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (env *Env) HealthCheck(c *gin.Context) {
	err := env.DB.DB.DB().Ping()
	if err != nil {

	}
	c.JSON(http.StatusOK, gin.H{"message": "healthy"})
}

func (env *Env) CreateToken(c *gin.Context) {
	token := env.J.NewAnonToken()

	c.JSON(http.StatusOK, gin.H{
		"expire": time.Now().Add(time.Duration(env.J.ExpiresAt) * time.Hour).Format(time.RFC3339),
		"token": token,
	})
}

func (env *Env) SigningKeyGet(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"alg": "RS512",
		"key": env.J.PublicString,
	})
}
