package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/mcavoyk/quirk/api/models"
	"github.com/sirupsen/logrus"
)

func (env *Env) SubmitVote(c *gin.Context) {
	vote := &models.Vote{}
	if err := c.ShouldBind(vote); err != nil {
		logrus.Debugf("Submit vote error: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Missing or invalid field 'vote'",
		})
		return
	}
	vote.UserID = c.GetString(UserContext)
	vote.PostID = c.Param("id")

	if err := env.WriteVote(vote); err != nil {
		logrus.Debugf("Submit vote error: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": err.Error()})
		return
	}
	c.JSON(http.StatusOK, http.StatusText(http.StatusOK))

}

func (env *Env) WriteVote(vote *models.Vote) error {
	if vote.Vote < -1 || vote.Vote > 1 {
		return fmt.Errorf("invalid vote request")
	}

	if err := env.db.Write(models.InsertVote, vote); err != nil {
		errNum := -1
		if sqlErr, ok := err.(*mysql.MySQLError); ok {
			errNum = int(sqlErr.Number)
		}
		if errNum == 1452 {
			return fmt.Errorf("post '%s' does not exist", vote.PostID)
		}
		return err
	}
	return nil
}
