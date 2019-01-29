package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/mcavoyk/quirk/api/location"
)

const Distance = 8.04672 // KM (5 Miles)

// Post represents top level content, viewable based on a user's location
// and the Posts Lat/Long
type Post struct {
	ID          string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	User        string `gorm:"index:user"`
	ParentID    string
	Depth       int     `gorm:"index:depth"`
	Content     string  `sql:"type:JSON"`
	Score       float64 `gorm:"-"`
	Positive    int     `gorm:"index:positive" json:"-" `
	Negative    int     `gorm:"index:negative" json:"-"`
	AccessType  string
	VoteState   int `gorm:"-"`
	NumComments int
	Collapsed   bool
	ColReason   string
	Lat         float64 `gorm:"index:latitude"`
	Lon         float64 `gorm:"index:longitude"`
}

func (db *DB) InsertPost(post *Post) (string, error) {
	if post.ParentID != "" {
		parent := db.GetPost(post.ParentID)
		if parent.ID == "" {
			return "", errors.New("invalid post parent")
		}
		post.Depth = parent.Depth + 1

	}
	post.ID = NewGUID()
	post.CreatedAt = time.Now()
	//db.Create(post)

	// Ignore error because this is a valid vote state
	_ = db.InsertOrUpdateVote(&Vote{User: post.User, PostID: post.ID, State: Upvote})
	return post.ID, nil
}

func (db *DB) GetPost(id string) *Post {
	post := new(Post)
	//db.Where("ID = ?", id).First(post)
	return post
}

func (db *DB) UpdatePost(post *Post) {
	post.ID = "" // Prevent user from updating primary key
	//db.Model(post).Updates(post)
	return
}

func (db *DB) DeletePost(id string) {
	if id == "" { // Gorm deletes all records if primary key is blank
		return
	}
	//db.Delete(&Post{ID: id})
	return
}

func (db *DB) PostsByDistance(lat, lon float64, page, pageSize int) []Post {
	posts := make([]Post, 0)

	points := location.BoundingPoints(&location.Point{lat, lon}, Distance)
	minLat := points[0].Lat
	minLon := points[0].Lon
	maxLat := points[1].Lat
	maxLon := points[1].Lon

	rows, err := db.Queryx("SELECT *, "+wilsonOrder+" FROM posts WHERE "+
		byDistance+" ORDER BY score DESC",
		minLat, maxLat, minLon, maxLon,
		lat, lat, lon, Distance/location.EarthRadius)

	if err != nil {
		fmt.Printf("SQL Error: %s\n", err.Error())
		return nil
	}

	defer rows.Close()

	start := (page - 1) * pageSize
	end := start + pageSize
	index := 0
	for true {
		if !rows.Next() || index > end {
			break
		}

		newPost := Post{}
		_ = rows.Scan(&newPost)
		posts = append(posts, newPost)
		index++
	}
	return posts
}

func (db *DB) PostsByParent(parentID string) []Post {
	posts := make([]Post, 0)

	rows, err := db.Queryx("WITH RECURSIVE cte as ("+
		"SELECT * FROM posts WHERE parent_id = ? UNION ALL "+
		"SELECT p.* FROM posts p INNER JOIN cte on p.parent_id = cte.id) "+
		"SELECT * FROM cte", parentID)
	if err != nil {
		fmt.Printf("SQL Error: %s\n", err.Error())
		return nil
	}

	defer rows.Close()
	for true {
		if !rows.Next() {
			break
		}

		newPost := Post{}
		err = rows.Scan(&newPost)
		posts = append(posts, newPost)
	}
	return posts

}

const wilsonOrder = "((positive + 1.9208) / (positive + negative) - 1.96 * SQRT((positive * negative) / (positive + negative)" +
	" + 0.9604) /(positive + negative)) / (1 + 3.8416 / (positive + negative)) AS score"

var byDistance = "(lat >= ? AND lat <= ?) AND (lon >= ? AND lon <= ?) AND ACOS(SIN(?) * SIN(lat) + COS(?) * COS(lat) * COS(lon - (?))) <= ?"
