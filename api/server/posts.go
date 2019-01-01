package server

import (
	"net/http"

	"../models"
	"github.com/gin-gonic/gin"
)

func (env *Env) HealthCheck(c *gin.Context) {
	err := env.db.DB.DB().Ping()
	if err != nil {

	}
	c.JSON(http.StatusOK, gin.H{"message": "healthy"})
}

func (env *Env) PostGet(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"Invalid": "Post ID can not be empty",
		})
		return
	}

	post := env.db.GetPost(id)
	c.JSON(http.StatusOK, post)
}

func (env *Env) PostDelete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"Invalid": "Post ID can not be empty",
		})
		return
	}

	env.db.DeletePost(id)
	c.JSON(http.StatusOK, http.StatusText(http.StatusOK))
}

func (env *Env) PostPatch(c *gin.Context) {
	post := &models.Post{}
	if err := c.Bind(post); err != nil {
		return
	}

	env.db.InsertPost(post)
	c.JSON(http.StatusOK, http.StatusText(http.StatusOK))
}

func (env *Env) PostPost(c *gin.Context) {
	post := &models.Post{}
	if err := c.Bind(post); err != nil {
		return
	}

	env.db.UpdatePost(post)
	c.JSON(http.StatusOK, http.StatusText(http.StatusOK))
}
