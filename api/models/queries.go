package models

const (
	SelectUserByName = "SELECT * FROM users WHERE username=?"
	InsertUser       = "INSERT INTO users (id, username, display_name, password, email)"
)
