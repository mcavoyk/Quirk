package models

import (
	"fmt"
)

// Vote represents a user's vote on a post
type Vote struct {
	User   string `gorm:"primary_key"`
	PostID string `gorm:"primary_key" binding:"required"`
	State  int    `binding:"required"`
}

const (
	Upvote   = 1
	Abstain  = 0
	Downvote = -1
)

// Valid vote states are -1, 0, 1; vote states of 0
// do not need to be stored as they represent no vote
func (db *DB) InsertOrUpdateVote(vote *Vote) error {
	if vote.State < Downvote || vote.State > Upvote {
		return fmt.Errorf("invalid vote state")
	}

	db.Exec("INSERT INTO votes (user, post_id, state) VALUES (?, ?, ?) "+
		"ON DUPLICATE KEY UPDATE state=VALUES(state)", vote.User, vote.PostID, vote.State)
	return nil
}

func (db *DB) GetVotesByUser(user string) []Vote {
	rows, err := db.Table("votes").Where("User = ?", user).Rows()
	if err != nil {
		fmt.Printf("SQL Error: %s\n", err.Error())
		return nil
	}

	defer rows.Close()

	votes := make([]Vote, 0)
	for true {
		if !rows.Next() {
			break
		}

		newVote := Vote{}
		_ = db.ScanRows(rows, &newVote)
		votes = append(votes, newVote)
	}
	return votes
}
