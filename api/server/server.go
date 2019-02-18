package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mcavoyk/quirk/api/models"
	"github.com/spf13/viper"
)

type Env struct {
	db models.Store
}

const (
	ApiV1 = "/api/v1"
	UserContext = "User"
	RootKey     = "RootKey"
)

func NewRouter(db models.Store, config *viper.Viper) http.Handler {
	levelStr := config.GetString("server.log_level")
	if levelStr != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	router.Use(setConfig(config))
	loadRoutes(router, &Env{db: db})
	return router
}

func setConfig(config *viper.Viper) gin.HandlerFunc {
	rootKey := config.GetString("server.root_key")
	if rootKey != "" {
		return func(c *gin.Context) {
			c.Set(RootKey, rootKey)
			c.Next()
		}
	}
	return func(c *gin.Context) {
		c.Next()
	}
}

func NewPass() string {
	return uuid.New().String()
}

func (env *Env) HasPermission(c *gin.Context, userID, resourceID string) error {
	if userID == resourceID {
		return nil
	}
	if val, ok := c.Get(RootKey); ok && val != nil {
		if val.(string) == userID {
			return nil
		}
	}

	return fmt.Errorf("Invalid permissions")
}

type Results struct {
	Page    int         `json:"page" form:"page,default=1" binding:"min=1"`
	PerPage int         `json:"per_page" form:"per_page,default=25" binding:"min=1"`
	Count   int         `json:"count"`
	Results interface{} `json:"results"`
}
