package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mcavoyk/quirk/api/models"
	"github.com/mcavoyk/quirk/api/pkg/location"
	"net/http"
)

type Post struct {
	Content    string  `json:"content" form:"content" binding:"required"`
	AccessType string  `json:"access_type" form:"access_type" binding:"required"`
	Lat        float64 `json:"lat" form:"lat" binding:"min=-90,max=90"`
	Lon        float64 `json:"lon" form:"lon" binding:"min=-180,max=180"`
}

func convertPost(src *Post, dst *models.Post) *models.Post {
	dst.Content = src.Content
	dst.AccessType = src.AccessType
	dst.Lat = location.ToRadians(src.Lat)
	dst.Lon = location.ToRadians(src.Lon)
	return dst
}


func (env *Env) GetPost(c *gin.Context) {
	id := c.Param("id")

	post, err := env.DB.GetPost(id, c.GetString(UserContext))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "Page not found"})
		return
	}
	c.JSON(http.StatusOK, post)
}

func (env *Env) DeletePost(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString(UserContext)

	post, err := env.DB.GetPost(id, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "Page not found"})
		return
	}

	if err := env.HasPermission(userID, post.UserID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"status": "Forbidden"})
		return
	}

	_ = env.DB.DeletePost(id)
	c.Status(http.StatusNoContent)
}

/*
func (env *Env) PatchPost(c *gin.Context) {
	id := c.Param("id")
	post := &Post{}
	if err := c.Bind(post); err != nil {
		return
	}
	existingPost := env.DB.GetPost(id)
	if existingPost.User != c.GetString(UserContext) {
		c.Status(http.StatusForbidden)
		return

	}
	newPost := convertPost(post, existingPost)
	env.DB.UpdatePost(newPost)
	c.JSON(http.StatusOK, http.StatusText(http.StatusOK))
}

*/
func (env *Env) CreatePost(c *gin.Context) {
	parentID := c.Param("id")
	givenPost := &Post{}
	if err := c.Bind(givenPost); err != nil {
		return
	}

	newPost := convertPost(givenPost, &models.Post{})
	newPost.UserID = c.GetString(UserContext)
	newPost.Parent = parentID

	post, err := env.DB.InsertPost(newPost)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, post)
}


// SearchPosts wraps search functions for posts
func (env *Env) SearchPosts(c *gin.Context) {
	coords, err := extractCoords(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": err.Error(),
		})
	}

	pageInfo := Results{}
	if err := c.ShouldBind(&pageInfo); err != nil {
		env.Log.Debugf("Search posts binding error: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Invalid fields for 'page' or 'per_page'",
		})
		return
	}

	posts, err := env.DB.PostsByDistance(coords.Lat, coords.Lon, c.GetString(UserContext), pageInfo.Page, pageInfo.PerPage, )

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": err.Error(),
		})
		return
	}

	pageInfo.Count = len(posts)
	pageInfo.Results = posts
	c.JSON(http.StatusOK, pageInfo)
}

func (env *Env) GetPostChildren(c *gin.Context) {
	pageInfo := Results{}
	parentID := c.Param("id")
	userID := c.GetString(UserContext)
	env.Log.Debugf("ParentID: %s | userID: %s", parentID, userID)

	if err := c.ShouldBind(&pageInfo); err != nil {
		env.Log.Debugf("Search posts binding error: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Invalid fields for 'page' or 'per_page'",
		})
		return
	}

	parentPost, err := env.DB.GetPost(parentID, userID)
	if err != nil {
		env.Log.Debugf("ParentID: %s | userID: %s", parentID, userID)
		c.JSON(http.StatusNotFound, gin.H{
			"status": "Page not found",
		})
		return
	}

	posts, err := env.DB.PostsByParent(fmt.Sprintf("%s/%s", parentPost.Parent, parentPost.ID), userID, pageInfo.Page, pageInfo.PerPage)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": err.Error(),
		})
		return
	}

	pageInfo.Count = len(posts)
	pageInfo.Results = posts
	c.JSON(http.StatusOK, pageInfo)
}