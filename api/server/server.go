package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mcavoyk/quirk/models"
)

type Env struct {
	DB    *models.DB
	Debug bool
}

const UserContext = "user"

func NewRouter(env *Env) http.Handler {
	if !env.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()
	loadRoutes(router, env)
	return router
}
