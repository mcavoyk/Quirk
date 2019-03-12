package models

import (
	"time"
)

// Session represents a User's login session
type Session struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Expiry    time.Time `json:"expiry"`
	IP        string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	Lat       float64   `json:"lat"`
	Lon       float64   `json:"lon"`
}

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

// PostInfo represents a specific user's "view" of a post
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

// Vote represents a user's vote on a post
type Vote struct {
	UserID string `json:"user_id" binding:"-"`
	PostID string `json:"post_id" binding:"-"`
	Vote   int    `json:"vote" form:"state" binding:"min=-1,max=1"`
}

type StoreBase interface {
	Close() error
	Ping() error
}
