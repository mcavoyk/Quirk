package server

import (
	"github.com/gin-gonic/gin"
)

func loadRoutes(router *gin.Engine, env *Env) {
	router.GET("/health", env.HealthCheck)

	if env.Auth {
		router.GET("/login", env.CreateToken)
		router.GET("/auth/signing", env.SigningKeyGet)
		router.Use(env.J.VerifyUser)
		router.GET("/auth/check", env.HealthCheck)
	}

	router.GET("/post/:id", env.PostGet)
	router.PATCH("/post/:id", env.PostPatch)
	router.DELETE("/post/:id", env.PostDelete)
	router.POST("/post", env.PostPost)

	router.POST("/vote", env.VotePost)

}
