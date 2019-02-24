package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/mcavoyk/quirk/api/models"
	"github.com/mcavoyk/quirk/api/pkg/location"
)

const PostDistance = 8.04672 // in KM (=5 Miles)

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

	newPost := marshalPost(givenPost)
	newPost.ID = models.NewGUID()
	newPost.UserID = c.GetString(UserContext)
	newPost.Parent = parentID
	if parentID != "" {
		var parentPost models.PostInfo
		err := env.db.ReadOne(&parentPost, models.SelectPostByUser, parentID, newPost.UserID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "Invalid parent post ID"})
			return
		}
	}

	err := env.db.Write(models.InsertPost, newPost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
		return
	}
	var post models.PostInfo
	err = env.db.ReadOne(&post, models.SelectPostByUser, newPost.ID, newPost.UserID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, post)
}

func (env *Env) GetPost(c *gin.Context) {
	id := c.Param("id")

	var post models.Post
	err := env.db.ReadOne(&post, models.SelectPostByUser, id, c.GetString(UserContext))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
		return
	}
	c.JSON(http.StatusOK, post)
}

func (env *Env) UpdatePost(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString(UserContext)
	post := &Post{}
	if err := c.Bind(post); err != nil {
		return
	}
	var existingPost models.Post
	err := env.db.ReadOne(&existingPost, models.SelectPostByUser, id, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "Page not found"})
		return
	}

	if err := env.HasPermission(c, userID, existingPost.UserID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"status": "Forbidden"})
		return
	}

	newPost := marshalPost(post)
	newPost.ID = id
	err = env.db.Write(models.UpdateValues("posts", *newPost), newPost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
		return
	}

	var updatedPost models.Post
	err = env.db.ReadOne(&updatedPost, models.SelectPostByUser, id, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedPost)
}

func (env *Env) DeletePost(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString(UserContext)

	var post models.Post
	err := env.db.ReadOne(&post, models.SelectPostByUser, id, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "Page not found"})
		return
	}

	if err := env.HasPermission(c, userID, post.UserID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"status": "Forbidden"})
		return
	}

	if _, err = env.db.Exec(models.DeletePostSoft, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
		return
	}
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

	posts := make([]models.PostInfo, 0)
	distArgs := byDistanceArgs(PostDistance, coords.Lat, coords.Lat)
	err = env.db.Read(&posts, models.SelectPostsByDistance, c.GetString(UserContext), distArgs[0], distArgs[1], distArgs[2], distArgs[3], distArgs[4], distArgs[5], distArgs[6], distArgs[7], pageInfo.PerPage, (pageInfo.Page-1)*pageInfo.PerPage)

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

	var parentPost models.Post
	err := env.db.ReadOne(&parentPost, models.SelectPostByUser, parentID, userID)
	if err != nil {
		logrus.Debugf("ParentID: %s | userID: %s", parentID, userID)
		c.JSON(http.StatusNotFound, gin.H{"status": "Page not found"})
		return
	}

	posts := make([]models.PostInfo, 0)
	err = env.db.Read(&posts, models.SelectPostsByParent, userID, fmt.Sprintf("%s/%s", parentPost.Parent, parentPost.ID), pageInfo.PerPage, (pageInfo.Page-1)*pageInfo.PerPage)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
		return
	}

	pageInfo.Count = len(posts)
	pageInfo.Results = posts
	c.JSON(http.StatusOK, pageInfo)
}

func marshalPost(post *Post) *models.Post {
	newPost := &models.Post{}
	bytes, err := json.Marshal(post.Content)
	if err != nil {
		fmt.Printf("Marshal error: %s\n", err.Error())
	}
	newPost.Content = string(bytes)
	newPost.AccessType = post.AccessType
	newPost.Lat = location.ToRadians(post.Lat)
	newPost.Lon = location.ToRadians(post.Lon)
	return newPost
}

func byDistanceArgs(distance, lat, lon float64) []float64 {
	lat, lon = location.ToRadians(lat), location.ToRadians(lon)

	points := location.BoundingPoints(&location.Point{Lat: lat, Lon: lon}, distance)
	minLat := points[0].Lat
	minLon := points[0].Lon
	maxLat := points[1].Lat
	maxLon := points[1].Lon

	logrus.Debugf("minLat %f | minLon %f | maxLat %f | maxLon %f | lat %f | lon %f", minLat, minLon, maxLat, maxLon, lat, lon)
	return []float64{minLat, maxLat, minLon, maxLon, lat, lat, lon, distance / location.EarthRadius}
}
