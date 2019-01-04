package server

import (
	"fmt"
	"github.com/mcavoyk/quirk/location"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mcavoyk/quirk/models"
)

type Post struct {
	Content    string
	AccessType string
	Lat        float64
	Lon        float64
}

func convertPost(src *Post, dst *models.Post) *models.Post {
	dst.Content = src.Content
	dst.AccessType = src.AccessType
	dst.Lat = location.ToRadians(src.Lat)
	dst.Lon = location.ToRadians(src.Lon)
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

func (env *Env) PostPost(c *gin.Context) {
	parentID := c.Param("postID")
	post := &Post{}
	if err := c.Bind(post); err != nil {
		return
	}

	newPost := convertPost(post, &models.Post{})
	newPost.User = c.GetString(UserContext)
	newPost.ParentID = parentID

	postID := env.DB.InsertPost(newPost)
	c.JSON(http.StatusOK, gin.H{
		"ID": postID,
		"ParentID": parentID,
	})
}

func (env *Env) PostsGet(c *gin.Context) {
	latStr := c.Query("lat")
	lonStr := c.Query("lon")
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pagesize", "25")

	lat, latErr := strconv.ParseFloat(latStr, 64)
	lon, lonErr := strconv.ParseFloat(lonStr, 64)

	if latErr != nil || lonErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Invalid latitude and longitude format",
		})
	}

	page, pageErr := strconv.Atoi(pageStr)
	pageSize, pageSizeErr := strconv.Atoi(pageSizeStr)

	if pageErr != nil || pageSizeErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Invalid page and pagesize format",
		})
		return
	}

	lat = location.ToRadians(lat)
	lon = location.ToRadians(lon)
	fmt.Printf("Received coords (%f, %f)\n", lat, lon)

	posts := env.DB.PostsByDistance(lat, lon, int(page), int(pageSize))
	c.JSON(http.StatusOK, posts)
}
