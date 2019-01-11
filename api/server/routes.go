package server

import (
	"github.com/gin-gonic/gin"
)

func loadRoutes(router *gin.Engine, env *Env) {

	api := router.Group("/api/v1")
	{
		api.GET("/health", env.HealthCheck)

		api.GET("/auth/token", env.CreateUser)
		api.GET("/auth/token/:token", env.ValidateUser)
		api.Use(env.UserVerify)
		router.NoRoute(noRoute)

		api.POST("/post", env.PostPost)
		api.GET("/post/:id", env.GetPost)
		api.PATCH("/post/:id", env.PatchPost)
		api.DELETE("/post/:id", env.DeletePost)

		api.POST("/post/:id/post", env.PostPost)
		api.GET("/post/:id/posts", env.GetPostsByPost)
		api.POST("/post/:id/vote", env.PostVote)

		api.GET("/posts", env.SearchPosts)
	}
}
