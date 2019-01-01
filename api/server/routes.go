package server

import (
	"../auth"
	"github.com/gin-gonic/gin"
)

func loadRoutes(router *gin.Engine, env *Env) {
	router.GET("/health", env.HealthCheck)
	router.Use(auth.VerifyUser)

	router.GET("/post/:id", env.PostGet)
	router.PATCH("/posts/:id", env.PostPatch)
	router.DELETE("/posts/:id", env.PostDelete)
	router.POST("/post", env.PostPost)

	router.GET("/comment/:id", env.CommentGet)
	router.PATCH("/comment/:id", env.CommentPatch)
	router.DELETE("/comment/:id", env.CommentDelete)
	router.POST("/comment", env.CommentPost)
}
