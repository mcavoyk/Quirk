package server

import (
	"net/http"

	"../models"
	"github.com/gin-gonic/gin"
)

func (env *Env) VotePost(c *gin.Context) {
	vote := &models.Vote{}
	if err := c.Bind(vote); err != nil {
		return
	}

	if err := env.db.InsertOrUpdateVote(vote); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, http.StatusText(http.StatusOK))
}
