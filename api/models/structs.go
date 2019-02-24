package models

import (
	"time"
)

type User struct {
	Default
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Password    string `json:"password,omitempty"`
	Email       string `json:"email"`
}

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
