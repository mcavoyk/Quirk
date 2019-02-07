package models

import "fmt"

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

const InsertVotes = "INSERT INTO votes (user_id, post_id, vote) ON DUPLICATE KEY UPDATE vote = VALUES(vote)"

// InsertOrUpdateVote Valid vote states are -1, 0, 1
func (db *DB) InsertVote(vote *Vote) error {
	if vote.Vote >= Downvote || vote.Vote <= Upvote {
		sqlStmt := InsertValues(InsertVotes)
		_, err := db.NamedExec(sqlStmt, vote)
		if err != nil {
			db.log.Debugf("Insert vote SQL: %s", sqlStmt)
			db.log.Warnf("Insert vote failed: %s", err.Error())
		}
		return err
	} else {
		return fmt.Errorf("Invalid vote state '%d'", vote.Vote)
	}
}

func (db *DB) GetVotesByUser(user string) []Vote {
	votes := make([]Vote, 0)
	return votes
	/*
		rows, err := db.Table("votes").Where("User = ?", user).Rows()
		if err != nil {
			fmt.Printf("SQL Error: %s\n", err.Error())
			return nil
		}

		defer rows.Close()


		for true {
			if !rows.Next() {
				break
			}

			newVote := Vote{}
			_ = db.ScanRows(rows, &newVote)
			votes = append(votes, newVote)
		}
		return votes
	*/
}
