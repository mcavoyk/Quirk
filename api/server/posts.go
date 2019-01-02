package server

import (
	"net/http"

	"../models"
	"github.com/gin-gonic/gin"
)

func (env *Env) PostGet(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"Invalid": "Post ID can not be empty",
		})
		return
	}

	post := env.DB.GetPost(id)
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

	env.DB.DeletePost(id)
	c.JSON(http.StatusOK, http.StatusText(http.StatusOK))
}

func (env *Env) PostPatch(c *gin.Context) {
	post := &models.Post{}
	if err := c.Bind(post); err != nil {
		return
	}

	env.DB.InsertPost(post)
	c.JSON(http.StatusOK, http.StatusText(http.StatusOK))
}

func (env *Env) PostPost(c *gin.Context) {
	post := &models.Post{}
	if err := c.Bind(post); err != nil {
		return
	}

	env.DB.UpdatePost(post)
	c.JSON(http.StatusOK, http.StatusText(http.StatusOK))
}
