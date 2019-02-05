package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mcavoyk/quirk/api/models"
	"github.com/mcavoyk/quirk/api/pkg/location"
)

type Post struct {
	Content    string  `json:"content" form:"content" binding:"required"`
	AccessType string  `json:"access_type" form:"access_type" binding:"required"`
	Lat        float64 `json:"lat" form:"lat" binding:"required"`
	Lon        float64 `json:"lon" form:"lon" binding:"required"`
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
func (env *Env) DeletePost(c *gin.Context) {
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

/*
// PostsGet wraps search functions for posts
func (env *Env) SearchPosts(c *gin.Context) {
	coords, err := extractCoords(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
	}

	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pagesize", "25")
	page, pageErr := strconv.Atoi(pageStr)
	pageSize, pageSizeErr := strconv.Atoi(pageSizeStr)

	if pageErr != nil || pageSizeErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Invalid page and pagesize format",
		})
		return
	}

	posts := env.DB.PostsByDistance(coords.Lat, coords.Lon, int(page), int(pageSize))
	votes := env.DB.GetVotesByUser(c.GetString(UserContext))
	env.Log.Debugf("Found %d votes submitted by user %s", len(votes), c.GetString(UserContext))

	//userCount := env.DB.UsersByDistance(coords.Lat, coords.Lon)
	//env.Log.Debugf("Found %d users in the radius of posts", userCount)


	for i := 0; i < len(posts); i++ {
		posts[i].Score = float64(posts[i].Positive - posts[i].Negative)
		for j := 0; j < len(votes); j++ {
			if posts[i].ID == votes[j].PostID {
				posts[i].VoteState = votes[j].State
				break
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"Posts": posts,
	})
}

func (env *Env) GetPostsByPost(c *gin.Context) {
	id := c.Param("id")
	posts := env.DB.PostsByParent(id)
	c.JSON(http.StatusOK, posts)
	return
}

// absoluteScore takes x - the wilson score of a post and the amount
// of users in the area and tries to create an absolute score
// using this to test a bit https://play.golang.org/p/iWDg0P9vXIP
func absoluteScore(x, totalVotes, users float64) float64 {
	//userScalar := math.Log(math.Max(float64(users), 1))
	startingWilsonScore := 0.206543

	// negative score
	shiftedScore := x - startingWilsonScore
	if (shiftedScore) < 0 {
		return math.Round(shiftedScore * totalVotes)
	}
	return math.Round(x * totalVotes)
}
*/
