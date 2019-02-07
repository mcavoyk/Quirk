package server

import (
	"github.com/gin-gonic/gin"
)

func loadRoutes(router *gin.Engine, env *Env) {

	api := router.Group("/api/v1")
	{
		api.GET("/health", env.HealthCheck)

		api.POST("/user", env.CreateUser)
		api.POST("/user/login", env.LoginUser)
		api.Use(env.UserVerify)
		router.NoRoute(noRoute)
		api.GET("/user/:id", env.GetUser)
		//api.PATCH("/user/:id", env.UpdateUser)
		api.DELETE("/user/:id", env.DeleteUser)

		api.POST("/post", env.CreatePost)
		api.GET("/post/:id", env.GetPost)
		api.PATCH("/post/:id", env.UpdatePost)
		api.DELETE("/post/:id", env.DeletePost)

		api.POST("/post/:id/post", env.CreatePost)
		api.GET("/post/:id/posts", env.GetPostChildren)
		api.POST("/post/:id/vote", env.SubmitVote)

		api.GET("/posts", env.SearchPosts)
	}
}
