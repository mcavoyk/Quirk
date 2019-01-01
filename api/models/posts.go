package models

import (
	"time"

	"github.com/segmentio/ksuid"
)

// LatLng represents a location on the Earth
type LatLng struct {
	Lat float64
	Lng float64
}

// Post represents a post on quirk
type Post struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time
	User      string
	Title     string
	Latitude  float64 `gorm:"index:latitude"`
	Longitude float64 `gorm:"index:longitude"`
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
	if id == "" { // Gorm deletes all records if primary key is blank,
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
