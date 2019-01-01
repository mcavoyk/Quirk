package server

import (
	"net/http"

	"../models"
	"github.com/gin-gonic/gin"
)

func (env *Env) CommentGet(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"Invalid": "Post ID can not be empty",
		})
		return
	}

	post := env.db.GetComment(id)
	c.JSON(http.StatusOK, post)
}

func (env *Env) CommentDelete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"Invalid": "Post ID can not be empty",
		})
		return
	}

	env.db.DeleteComment(id)
	c.JSON(http.StatusOK, http.StatusText(http.StatusOK))
}

func (env *Env) CommentPatch(c *gin.Context) {
	comment := &models.Comment{}
	if err := c.Bind(comment); err != nil {
		return
	}

	env.db.InsertComment(comment)
	c.JSON(http.StatusOK, http.StatusText(http.StatusOK))
}

func (env *Env) CommentPost(c *gin.Context) {
	comment := &models.Comment{}
	if err := c.Bind(comment); err != nil {
		return
	}

	env.db.UpdateComment(comment)
	c.JSON(http.StatusOK, http.StatusText(http.StatusOK))
}
