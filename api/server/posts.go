package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/mcavoyk/quirk/api/models"
	"github.com/mcavoyk/quirk/api/pkg/location"
)

type Post struct {
	Content    Content `json:"content" form:"content" binding:"required"`
	AccessType string  `json:"access_type" form:"access_type" binding:"required"`
	Lat        float64 `json:"lat" form:"lat" binding:"min=-90,max=90"`
	Lon        float64 `json:"lon" form:"lon" binding:"min=-180,max=180"`
}

type Content struct {
	Title string `json:"title" binding:"required"`
}

func (env *Env) CreatePost(c *gin.Context) {
	parentID := c.Param("id")
	givenPost := new(Post)
	if err := c.ShouldBind(givenPost); err != nil {
		logrus.Debugf("Create posts binding error: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "Invalid request"})
		return
	}

	newPost := convertPost(givenPost, &models.Post{})
	newPost.UserID = c.GetString(UserContext)
	newPost.Parent = parentID

	post, err := env.DB.InsertPost(newPost)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": err.Error()})
		return
	}

	c.JSON(http.StatusOK, post)
}

func (env *Env) GetPost(c *gin.Context) {
	start := time.Now()
	id := c.Param("id")

	post, err := env.DB.GetPostByUser(id, c.GetString(UserContext))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "Page not found"})
		return
	}
	c.JSON(http.StatusOK, post)
	logrus.Debugf("Total get post time: %f", time.Since(start).Seconds())
}

func (env *Env) UpdatePost(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString(UserContext)
	post := &Post{}
	if err := c.Bind(post); err != nil {
		return
	}
	existingPost, err := env.DB.GetPostByUser(id, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "Page not found"})
		return
	}

	if err := env.HasPermission(c, userID, existingPost.UserID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"status": "Forbidden"})
		return
	}

	newPost := convertPost(post, &models.Post{})
	newPost.ID = id
	returnedPost, err := env.DB.UpdatePost(newPost, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
		return
	}
	c.JSON(http.StatusOK, returnedPost)
}

func (env *Env) DeletePost(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString(UserContext)

	post, err := env.DB.GetPostByUser(id, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "Page not found"})
		return
	}

	if err := env.HasPermission(c, userID, post.UserID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"status": "Forbidden"})
		return
	}

	_ = env.DB.DeletePost(id)
	c.Status(http.StatusNoContent)
}

// SearchPosts wraps search functions for posts
func (env *Env) SearchPosts(c *gin.Context) {
	coords, err := extractCoords(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": err.Error()})
	}

	pageInfo := Results{}
	if err := c.ShouldBind(&pageInfo); err != nil {
		logrus.Debugf("Search posts binding error: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "Invalid fields for 'page' or 'per_page'"})
		return
	}

	posts, err := env.DB.PostsByDistance(coords.Lat, coords.Lon, c.GetString(UserContext), pageInfo.Page, pageInfo.PerPage)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
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

	if err := c.ShouldBind(&pageInfo); err != nil {
		logrus.Debugf("Search posts binding error: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "Invalid fields for 'page' or 'per_page'"})
		return
	}

	parentPost, err := env.DB.GetPostByUser(parentID, userID)
	if err != nil {
		logrus.Debugf("ParentID: %s | userID: %s", parentID, userID)
		c.JSON(http.StatusNotFound, gin.H{"status": "Page not found"})
		return
	}

	posts, err := env.DB.PostsByParent(fmt.Sprintf("%s/%s", parentPost.Parent, parentPost.ID), userID, pageInfo.Page, pageInfo.PerPage)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
		return
	}

	pageInfo.Count = len(posts)
	pageInfo.Results = posts
	c.JSON(http.StatusOK, pageInfo)
}

func convertPost(src *Post, dst *models.Post) *models.Post {
	bytes, err := json.Marshal(src.Content)
	if err != nil {
		fmt.Printf("Marshal error: %s\n", err.Error())
	}
	dst.Content = string(bytes)
	dst.AccessType = src.AccessType
	dst.Lat = location.ToRadians(src.Lat)
	dst.Lon = location.ToRadians(src.Lon)
	return dst
}
