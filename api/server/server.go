package server

import (
	"github.com/gin-gonic/gin"
	"github.com/mcavoyk/quirk/api/models"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Env struct {
	DB    *models.DB
	Log   *logrus.Logger
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
