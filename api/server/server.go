package server

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mcavoyk/quirk/api/models"
)

type Env struct {
	DB    *models.DB
	Debug bool
	log *log.Logger
}

const UserContext = "user"

func NewRouter(env *Env) http.Handler {
	env.log = log.New(os.Stdout, "", log.Ltime)
	if !env.Debug {
		gin.SetMode(gin.ReleaseMode)
		env.log = log.New(os.Stdout, "", log.Ltime)

	}
	router := gin.Default()
	loadRoutes(router, env)
	return router
}
