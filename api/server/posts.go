package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mcavoyk/quirk/models"
)

type Post struct {
	ParentID   string
	Title      string
	Content    interface{}
	AccessType string
	Latitude   float64
	Longitude  float64
}

func convertPost(src *Post) *models.Post {
	return &models.Post{}
}

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

	currentPost := env.DB.GetPost(id)
	if currentPost.User != c.GetString(UserContext) {
		c.JSON(http.StatusForbidden, gin.H{
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

	env.DB.UpdatePost(post)
	c.JSON(http.StatusOK, http.StatusText(http.StatusOK))
}

func (env *Env) PostPost(c *gin.Context) {
	post := &Post{}
	if err := c.Bind(post); err != nil {
		return
	}

	newPost := convertPost(post)
	newPost.User = c.GetString(UserContext)

	env.DB.InsertPost(newPost)
	c.JSON(http.StatusOK, http.StatusText(http.StatusOK))
}
