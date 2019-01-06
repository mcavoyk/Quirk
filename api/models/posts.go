package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/mcavoyk/quirk/api/location"
)

// Post represents top level content, viewable based on a user's location
// and the Posts Lat/Long
type Post struct {
	ID          string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	User        string `gorm:"index:user"`
	ParentID    string
	Depth       int    `gorm:"index:depth"`
	Content     string `sql:"type:JSON"`
	Score       int    `gorm:"index:score"`
	AccessType  string
	Vote        []Vote `gorm:"ForeignKey:ID"`
	VoteState   int    `gorm:"-"`
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
	db.Create(post)
	return post.ID, nil
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

func (db *DB) PostsByDistance(lat, lon float64, page, pageSize int) []Post {
	posts := make([]Post, 0)
	distance := 5.0 // KM
	points := location.BoundingPoints(&location.Point{lat, lon}, distance)
	minLat := points[0].Lat
	minLon := points[0].Lon
	maxLat := points[1].Lat
	maxLon := points[1].Lon

	fmt.Printf("Using (%f, %f), calculated bounding points:\n (%f, %f) and (%f, %f)\n", lat, lon, minLat, minLon, maxLat, maxLon)

	sql := fmt.Sprintf("SELECT * FROM posts WHERE "+
		"(lat >= %f AND lat <= %f) AND (lon >= %f AND lon <= %f) "+
		"AND ACOS(SIN(%f) * SIN(lat) + COS(%f) * COS(lat) * COS(lon - (%f))) <= %f",
		minLat, maxLat, minLon, maxLon,
		lat, lat, lon, distance/location.EarthRadius)

	fmt.Printf("%s\n", sql)
	rows, err := db.Raw(sql).Rows()
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
		_ = db.ScanRows(rows, &newPost)
		posts = append(posts, newPost)
		index++
	}
	return posts
}

func (db *DB) PostsByParent(parentID string) []Post {
	posts := make([]Post, 0)

	sql := fmt.Sprintf("WITH RECURSIVE cte (id, parent_id) as ("+
		"SELECT id, parent_id FROM posts WHERE parent_id = '%s' UNION ALL "+
		"SELECT p.id, p.parent_id FROM posts p INNER JOIN cte on p.parent_id = cte.id) "+
		"SELECT * FROM cte", parentID)
	fmt.Printf("%s\n", sql)

	rows, err := db.Raw(sql).Rows()
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
		_ = db.ScanRows(rows, &newPost)
		posts = append(posts, newPost)
	}
	return posts

}
