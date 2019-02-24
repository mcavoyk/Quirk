package server

import (
	"github.com/gin-gonic/gin"
)

func loadRoutes(router *gin.Engine, env *Env) {
	router.NoRoute(noRoute)
	api := router.Group(ApiV1)
	{
		api.GET("/health", env.healthCheck)

		api.POST("/user", env.CreateUser)
		api.POST("/user/login", env.LoginUser)
		api.Use(env.UserVerify)
		//api.GET("/db", env.selectQuery)
		api.GET("/user/:id", env.GetUser)
		api.PATCH("/user/:id", env.PatchUser)
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
