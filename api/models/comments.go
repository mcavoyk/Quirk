package models

import (
	"time"

	"github.com/segmentio/ksuid"
)

// Comment represents a comment on quirk
type Comment struct {
	ID        string
	ParentID  string `gorm:"index:parent"`
	CreatedAt time.Time
	UpdatedAt time.Time
	User      string
	Contents  string
}

func (db *DB) InsertComment(comment *Comment) {
	comment.ID = ksuid.New().String()
	db.Create(comment)
	return
}

func (db *DB) GetComment(id string) *Comment {
	comment := new(Comment)
	db.Where("ID = ?", id).First(comment)
	return comment
}

func (db *DB) UpdateComment(comment *Comment) {
	comment.ID = "" // Prevent user from updating primary key
	db.Model(comment).Updates(comment)
	return
}

func (db *DB) DeleteComment(id string) {
	if id == "" { // Gorm deletes all records if primary key is blank,
		return
	}
	db.Delete(&Comment{ID: id})
	return
}
