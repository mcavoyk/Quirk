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

// Valid vote states are -1, 0, 1; vote states of 0
// do not need to be stored as they represent no vote
func (db *DB) InsertOrUpdateVote(vote *Vote) error {
	if vote.State < -1 || vote.State > 1 {
		return fmt.Errorf("invalid vote state")
	}

	sql := fmt.Sprintf("INSERT INTO votes (user, post_id, state) VALUES ('%s', '%s', %d) "+
		"ON DUPLICATE KEY UPDATE state=VALUES(state)",
		vote.User, vote.PostID, vote.State)
	db.Exec(sql)
	return nil
}

func (db *DB) GetVotesByUser(user string) *Vote {

}
