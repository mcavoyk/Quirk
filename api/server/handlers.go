package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (env *Env) HealthCheck(c *gin.Context) {
	err := env.DB.DB.DB().Ping()
	if err != nil {

	}
	c.JSON(http.StatusOK, gin.H{"message": "healthy"})
}

func noRoute(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"message": "Page not found",
	})
}
