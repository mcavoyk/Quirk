package models

// Vote represents a user's vote on a post
type Vote struct {
	UserID string `json:"user_id" binding:"-"`
	PostID string `json:"post_id" binding:"-"`
	Vote   int    `json:"vote" form:"state" binding:"min=-1,max=1"`
}

const (
	Upvote   = 1
	Downvote = -1
)

// InsertOrUpdateVote Valid vote states are -1, 0, 1
/*
func (db *DB) InsertVote(vote *Vote) error {
	if vote.Vote >= Downvote || vote.Vote <= Upvote {
		sqlStmt := InsertValues(InsertVotes)
		_, err := db.NamedExec(sqlStmt, vote)
		if err != nil {
			logrus.Debugf("Insert vote SQL: %s", sqlStmt)
			logrus.Warnf("Insert vote failed: %s", err.Error())
		}
		return err
	} else {
		return fmt.Errorf("Invalid vote state '%d'", vote.Vote)
	}
}
*/
