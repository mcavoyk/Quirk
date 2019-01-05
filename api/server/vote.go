package server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mcavoyk/quirk/models"
)

func (env *Env) PostVote(c *gin.Context) {
	postID := c.Param("id")
	user := c.GetString(UserContext)
	stateStr := c.Query("state")

	state, err := strconv.Atoi(stateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Invalid or missing format for value 'state'",
		})
		return
	}

	newVote := &models.Vote{
		PostID: postID,
		User:   user,
		State:  state,
	}
	env.log.Printf("Received vote: %+v\n", newVote)
	if err := env.DB.InsertOrUpdateVote(newVote); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, http.StatusText(http.StatusOK))
}
