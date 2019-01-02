package server

import (
	"../auth"
	"github.com/gin-gonic/gin"
)

func loadRoutes(router *gin.Engine, env *Env) {
	router.GET("/health", env.HealthCheck)
	router.GET("/login", env.CreateToken)
	router.GET("/auth/signing", env.SigningKeyGet)
	router.Use(auth.VerifyUser)

	router.GET("/post/:id", env.PostGet)
	router.PATCH("/post/:id", env.PostPatch)
	router.DELETE("/post/:id", env.PostDelete)
	router.POST("/post", env.PostPost)

	router.POST("/vote", env.VotePost)

}
