package models

// Vote represents a user's vote on a post
type Vote struct {
	User   string `db:"user_id"`
	PostID string `binding:"required"`
	State  int    `binding:"required"`
}

const (
	Upvote   = 1
	Abstain  = 0
	Downvote = -1
)

// InsertOrUpdateVote Valid vote states are -1, 0, 1
func (db *DB) InsertOrUpdateVote(vote *Vote) error {
	_, err := db.NamedExec("INSERT INTO votes (user_id, post_id, state) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE state=VALUES(state)", vote)
	return err
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

