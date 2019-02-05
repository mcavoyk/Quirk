package models

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mcavoyk/quirk/api/pkg/location"
)

const Distance = 8.04672 // KM (5 Miles)

// Post represents top level content, viewable based on a user's location
// and the Posts Lat/Long
type Post struct {
	Default
	UserID     string  `json:"user_id"`
	Parent     string  `json:"parent"`
	Lat        float64 `json:"lat"`
	Lon        float64 `json:"lon"`
	AccessType string  `json:"access_type"`
	Content    string  `json:"content"`
}

type PostInfo struct {
	Post
	Positive    int     `json:"positive"`
	Negative    int     `json:"negative"`
	Score       float64 `json:"score"`
	Username    string  `json:"username"`
	DisplayName string  `json:"display_name"`
	VoteState   int     `json:"vote_state"`
	NumChildren int     `json:"num_children"`
	//Collapsed   bool
	//ColReason   string
}

const InsertPost = "INSERT INTO posts (id, user_id, parent, lat, lon, access_type, content)"

func (db *DB) InsertPost(post *Post) (*PostInfo, error) {
	if post.Parent != "" {
		parentSplit := strings.Split(post.Parent, "/")
		lastParent := parentSplit[len(parentSplit)-1]
		parent, err := db.GetPost(lastParent, post.UserID)
		if err != nil {
			return nil, errors.New("invalid post parent")
		}
		post.Parent = fmt.Sprintf("%s/%s", parent.Parent, lastParent)
	}

	post.ID = NewGUID()
	sqlStmt := InsertValues(InsertPost)
	db.log.Debugf("Insert post statement: %s", sqlStmt)
	_, err := db.NamedExec(sqlStmt, post)
	if err != nil {
		db.log.Warnf("Insert post failed: %s", err.Error())
		return nil, err
	}

	_ = db.InsertVote(&Vote{UserID: post.UserID, PostID: post.ID, Vote: Upvote})
	return db.GetPost(post.ID, post.UserID)
}

func (db *DB) GetPost(id string, user string) (*PostInfo, error) {
	post := new(PostInfo)
	err := db.Unsafe().Get(post, "SELECT * FROM post_view WHERE id=? AND vote_user_id=?", id, user)
	if err != nil {
		db.log.Debugf("Get post failed: %s", err.Error())
		return nil, err
	}
	post.Lat, post.Lon = location.ToDegrees(post.Lat), location.ToDegrees(post.Lon)
	return post, nil
}

func (db *DB) UpdatePost(post *Post) {
	post.ID = "" // Prevent user from updating primary key
	//db.Model(post).Updates(post)
	return
}

func (db *DB) DeletePost(id string) error {
	_, err := db.Exec("UPDATE posts SET deleted_at = NOW() WHERE id=?", id)
	if err != nil {
		db.log.Debugf("Delete post failed: %s", err.Error())
	}
	return err
}

func (db *DB) PostsByDistance(lat, lon float64, userID string, page, pageSize int) ([]PostInfo, error) {
	lat, lon = location.ToRadians(lat), location.ToRadians(lon)
	posts := []PostInfo{}

	points := location.BoundingPoints(&location.Point{lat, lon}, Distance)
	minLat := points[0].Lat
	minLon := points[0].Lon
	maxLat := points[1].Lat
	maxLon := points[1].Lon

	db.log.Debugf("minLat %f | minLon %f | maxLat %f | maxLon %f | lat %f | lon %f", minLat,  minLon, maxLat, maxLon, lat, lon)
	err := db.Unsafe().Select(&posts, "SELECT * FROM post_view WHERE deleted_at IS NULL AND vote_user_id = ? AND "+
		byDistance+" ORDER BY score DESC LIMIT ? OFFSET ?",
		userID, minLat, maxLat, minLon, maxLon, lat, lat, lon, Distance/location.EarthRadius,
		pageSize, (page - 1) * pageSize)

	if err != nil {
		db.log.Errorf("Select posts by distance error: %s", err.Error())
		return nil, err
	}

	return posts, nil
}

func (db *DB) PostsByParent(parent, user string, page, pageSize int) ([]PostInfo, error) {
	posts := make([]PostInfo, 0)

	err := db.Unsafe().Select(&posts,"SELECT * FROM post_view WHERE vote_user_id = ? AND " +
		"parent LIKE CONCAT(?, '%') ORDER BY score DESC LIMIT ? OFFSET ?", user, parent, pageSize, (page - 1) * pageSize)

	if err != nil {
		db.log.Errorf("Select posts by parent error: %s", err.Error())
		return nil, err
	}

	return posts, nil
}

const wilsonOrder = "((positive + 1.9208) / (positive + negative) - 1.96 * SQRT((positive * negative) / (positive + negative)" +
	" + 0.9604) /(positive + negative)) / (1 + 3.8416 / (positive + negative)) AS score"

var byDistance = "(lat >= ? AND lat <= ?) AND (lon >= ? AND lon <= ?) AND ACOS(SIN(?) * SIN(lat) + COS(?) * COS(lat) * COS(lon - (?))) <= ?"
