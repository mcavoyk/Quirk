package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mcavoyk/quirk/models"
)

type Post struct {
	ID         string
	ParentID   string
	Content    string
	AccessType string
	Latitude   float64
	Longitude  float64
}

func convertPost(src *Post, dst *models.Post) *models.Post {
	dst.ParentID = src.ParentID
	dst.Content = src.Content
	dst.AccessType = src.AccessType
	dst.Latitude = src.Latitude
	dst.Longitude = src.Longitude
	return dst
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
		c.Status(http.StatusForbidden)
		return
	}
	env.DB.DeletePost(id)
	c.JSON(http.StatusOK, http.StatusText(http.StatusOK))
}

func (env *Env) PostPatch(c *gin.Context) {
	post := &Post{}
	if err := c.Bind(post); err != nil {
		return
	}
	existingPost := env.DB.GetPost(post.ID)
	if existingPost.User != c.GetString(UserContext) {
		c.Status(http.StatusForbidden)
		return

	}
	newPost := convertPost(post, existingPost)
	env.DB.UpdatePost(newPost)
	c.JSON(http.StatusOK, http.StatusText(http.StatusOK))
}

func (env *Env) PostPost(c *gin.Context) {
	post := &Post{}
	if err := c.Bind(post); err != nil {
		return
	}

	newPost := convertPost(post, &models.Post{})
	newPost.User = c.GetString(UserContext)

	postID := env.DB.InsertPost(newPost)
	c.JSON(http.StatusOK, gin.H{
		"ID": postID,
	})
}
