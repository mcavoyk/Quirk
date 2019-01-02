package server

import (
	"net/http"

	"../auth"
	"../models"
	"github.com/gin-gonic/gin"
)

type Env struct {
	DB *models.DB
	J  *auth.JWTStorage
}

func NewRouter(env *Env) http.Handler {
	router := gin.Default()
	loadRoutes(router, env)
	return router
}
