package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mcavoyk/quirk/api/models"
	"github.com/sirupsen/logrus"
)

type Env struct {
	DB  *models.DB
	Log *logrus.Logger
}

const UserContext = "user"

func NewRouter(env *Env) http.Handler {
	if env.Log.Level != logrus.DebugLevel {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	loadRoutes(router, env)
	return router
}

func NewPass() string {
	return uuid.New().String()
}

func (env *Env) HasPermission(userID, resourceID string) error {
	if userID == resourceID {
		return nil
	}
	return fmt.Errorf("Invalid permissions")
}