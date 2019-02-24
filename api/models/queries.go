package models

const (
	SelectUser          = "SELECT * FROM users WHERE id=?"
	SelectUserByName    = "SELECT * FROM users WHERE username=?"
	SelectSession       = "SELECT * FROM sessions WHERE id=?"
	SelectUserBySession = "SELECT * FROM user_sessions WHERE session_id=?"
	UpdateSession       = "UPDATE sessions SET ip_address = :ip_address, lat = :lat, lon = :lon WHERE id = :id"
	DeleteUserSoft      = "UPDATE users SET deleted_at = NOW() WHERE id = ?"
	DeleteSessions      = "DELETE FROM sessions WHERE user_id = ?"

	SelectPostByUser      = "SELECT * FROM post_view WHERE id=? AND vote_user_id=?"
	DeletePostSoft        = "UPDATE posts SET deleted_at = NOW() WHERE id=?"
	SelectPostsByDistance = "SELECT * FROM post_view WHERE deleted_at IS NULL AND vote_user_id = ? AND " + byDistance + " ORDER BY score DESC LIMIT ? OFFSET ?"
	SelectPostsByParent   = "SELECT * FROM post_view WHERE vote_user_id = ? AND parent LIKE CONCAT(?, '%') ORDER BY score DESC LIMIT ? OFFSET ?"

	// helper query parts
	byDistance  = "(lat >= ? AND lat <= ?) AND (lon >= ? AND lon <= ?) AND ACOS(SIN(?) * SIN(lat) + COS(?) * COS(lat) * COS(lon - (?))) <= ?"
	wilsonOrder = "((positive + 1.9208) / (positive + negative) - 1.96 * SQRT((positive * negative) / (positive + negative)" +
		" + 0.9604) /(positive + negative)) / (1 + 3.8416 / (positive + negative)) AS score"
)

var (
	InsertUser    = InsertValues("INSERT INTO users (id, username, display_name, password, email)")
	InsertSession = InsertValues("INSERT INTO sessions (id, user_id, expiry, ip_address, user_agent, lat, lon)")
	InsertPost    = InsertValues("INSERT INTO posts (id, user_id, parent, lat, lon, access_type, content)")
	InsertVote    = InsertValues("INSERT INTO votes (user_id, post_id, vote) ON DUPLICATE KEY UPDATE vote = VALUES(vote)")
)
