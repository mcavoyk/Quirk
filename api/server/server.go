package server

import (
	"net/http"

	"../models"
	"github.com/gin-gonic/gin"
)

type Env struct {
	db *models.DB
}

func NewRouter(db *models.DB) http.Handler {
	router := gin.Default()
	env := &Env{db: db}
	loadRoutes(router, env)
	return router
}
