package models

import (
	"time"

	"github.com/segmentio/ksuid"
)

// Post represents top level content, viewable based on a user's location
// and the Posts Lat/Long
type Post struct {
	ID          string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	User        string `gorm:"index:user"`
	ParentID    string
	Depth       int `gorm:"index:depth"`
	Title       string
	Body        string `sql:"type:text"`
	Score       int    `gorm:"index:score"`
	AccessType  string
	Vote        []Vote `gorm:"ForeignKey:ID"`
	VoteState   int    `gorm:"-"`
	NumComments int
	Collapsed   bool
	ColReason   string
	Latitude    float64 `gorm:"index:latitude"`
	Longitude   float64 `gorm:"index:longitude"`
}

func (db *DB) InsertPost(post *Post) {
	post.ID = ksuid.New().String()
	db.Create(post)
	return
}

func (db *DB) GetPost(id string) *Post {
	post := new(Post)
	db.Where("ID = ?", id).First(post)
	return post
}

func (db *DB) UpdatePost(post *Post) {
	post.ID = "" // Prevent user from updating primary key
	db.Model(post).Updates(post)
	return
}

func (db *DB) DeletePost(id string) {
	if id == "" { // Gorm deletes all records if primary key is blank
		return
	}
	db.Delete(&Post{ID: id})
	return
}

func (db *DB) GetPosts() []Post {
	posts := make([]Post, 0)
	db.Find(&posts)
	return posts
}
