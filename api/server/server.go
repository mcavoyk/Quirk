package server

import (
	"net/http"

	"../auth"
	"../models"
	"github.com/gin-gonic/gin"
)

type Env struct {
	DB    *models.DB
	J     *auth.JWTStorage
	Auth  bool
	Debug bool
}

func NewRouter(env *Env) http.Handler {
	if !env.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()
	loadRoutes(router, env)
	return router
}
