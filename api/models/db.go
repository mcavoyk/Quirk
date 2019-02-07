package models

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx/reflectx"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	_ "github.com/go-sql-driver/mysql"
	"github.com/segmentio/ksuid"
)

type DB struct {
	*sqlx.DB
	log *logrus.Logger
}

type Default struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt *NullTime `json:"deleted_at,omitempty"`
}

// InitDB
func InitDB(connection string) (*DB, error) {
	connectionParams := "?parseTime=True&charset=utf8mb4&collation=utf8mb4_unicode_ci"
	db, err := sqlx.Connect("mysql", connection+connectionParams)
	if err != nil {
		return nil, err
	}

	for _, stmt := range schema {
		_, err = db.Exec(stmt)
		if err != nil {
			return nil, err
		}
	}

	db, err = sqlx.Connect("mysql", connection+dbName+connectionParams)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Minute * 10)
	db.Mapper = reflectx.NewMapperFunc("json", strings.ToLower)

	return &DB{DB: db, log: logrus.New()}, nil
}

func NewGUID() string {
	return ksuid.New().String()
}

func (db *DB) SetLogLevel(logLevel string) {
	if level, err := logrus.ParseLevel(logLevel); err == nil {
		db.log.SetLevel(level)
	}
}

func InsertValues(insert string) string {
	start := strings.Split(insert, "(")
	if len(start) == 1 {
		return insert
	}
	start[1] = strings.Join(start[1:], "(")

	end := strings.Split(start[1], ")")
	endStmt := ""
	if len(end) > 1 {
		endStmt = strings.Join(end[1:], ")")
	}

	columns := end[0]
	columnSplit := strings.Split(columns, ",")
	namedValues := ""
	for i, column := range columnSplit {
		column = strings.TrimSpace(column)
		if i != 0 {
			namedValues += ", "
		}
		namedValues += fmt.Sprintf(":%s", column)
	}

	return fmt.Sprintf("%s(%s) VALUES (%s)%s", start[0], end[0], namedValues, endStmt)

}

// NullTime represents a time.Time that may be null. NullTime implements the
// sql.Scanner interface so it can be used as a scan destination, similar to
// sql.NullString.
type NullTime struct {
	time.Time
	Null bool // Valid is true if Time is not NULL
}

// Scan implements the Scanner interface.
func (nt *NullTime) Scan(value interface{}) error {
	nt.Time, nt.Null = value.(time.Time)
	nt.Null = !nt.Null
	return nil
}

// Value implements the driver Valuer interface.
func (nt NullTime) Value() (driver.Value, error) {
	if nt.Null {
		return nil, nil
	}
	return nt.Time, nil
}

const dbName = "quirk"

// Database schema
var schema = []string{
	`CREATE DATABASE IF NOT EXISTS ` + dbName + ` CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci;`,

	`USE quirk;`,
	`CREATE TABLE IF NOT EXISTS metadata (
  		name VARCHAR(255),
		version VARCHAR(255),
  		updated_at TIMESTAMP NOT NULL DEFAULT NOW() ON UPDATE NOW(),
  		PRIMARY KEY (name)
	);`,

	`CREATE TABLE IF NOT EXISTS users (
 		id VARCHAR(255) PRIMARY KEY NOT NULL,
  		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW() ON UPDATE NOW(),
  		deleted_at TIMESTAMP NULL,
  		username VARCHAR(255) UNIQUE NOT NULL,
		display_name VARCHAR(255) NOT NULL,
		password VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL DEFAULT '',
 
		INDEX idx_username (username),
		INDEX idx_email (email),
		INDEX idx_created (created_at),
		INDEX idx_updated (updated_at),
		INDEX idx_deleted (deleted_at)
	);`,

	`CREATE TABLE IF NOT EXISTS sessions (
 		id VARCHAR(255) PRIMARY KEY NOT NULL,
		user_id VARCHAR(255) NOT NULL,
  		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW() ON UPDATE NOW(),
		expiry TIMESTAMP NOT NULL,
		ip_address VARCHAR(255) NOT NULL DEFAULT '',
		user_agent VARCHAR(255) NOT NULL DEFAULT '',
		lat DOUBLE NOT NULL,
  		lon DOUBLE NOT NULL,

		INDEX idx_latlon (lat, lon),
		INDEX idx_created (created_at),
		INDEX idx_updated (updated_at),
		INDEX idx_expiry (expiry),

		FOREIGN KEY fk_user_sessions (user_id)
			REFERENCES users(id)
			ON DELETE CASCADE
	);`,

	`CREATE TABLE IF NOT EXISTS posts (
  		id VARCHAR(255) PRIMARY KEY NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW() ON UPDATE NOW(),
  		deleted_at TIMESTAMP NULL,
		lat DOUBLE NOT NULL,
		lon DOUBLE NOT NULL,
		user_id VARCHAR(255),
		parent VARCHAR(1023) NOT NULL DEFAULT '',
		access_type ENUM('public', 'private') NOT NULL DEFAULT 'public',
		content TEXT NOT NULL,

		INDEX idx_parent (parent),
		INDEX idx_latlon (lat, lon),
		INDEX idx_created (created_at),
		INDEX idx_deleted (updated_at),

		FOREIGN KEY fk_user_posts (user_id)
			REFERENCES users(id)
			ON DELETE SET NULL
	);`,

	`CREATE TABLE IF NOT EXISTS votes (
  		post_id VARCHAR(255) NOT NULL,
  		user_id VARCHAR(255) NOT NULL,
  		vote TINYINT NOT NULL,
  		updated_at TIMESTAMP NOT NULL DEFAULT NOW() ON UPDATE NOW(),
  		PRIMARY KEY (post_id, user_id),

  		INDEX idx_post (post_id),
  		INDEX idx_user (user_id),
		INDEX idx_vote (vote),

  		FOREIGN KEY fk_user_votes (user_id)
    		REFERENCES users(id)
    		ON DELETE SET NULL,

  		FOREIGN KEY fk_post_votes (post_id)
  		  REFERENCES posts(id)
    		ON DELETE CASCADE
	);`,

	`CREATE VIEW IF NOT EXISTS posts_live AS
		SELECT id, created_at, updated_at, deleted_at, lat, lon, parent, access_type, user_id,
   	 	(CASE WHEN deleted_at IS NOT NULL THEN '[deleted]' ELSE content END) AS content
    FROM posts`,

    `CREATE VIEW IF NOT EXISTS users_live AS
		SELECT id, created_at, updated_at, deleted_at,
    	(CASE WHEN deleted_at IS NOT NULL THEN '[deleted]' ELSE username END) AS username,
    	(CASE WHEN deleted_at IS NOT NULL THEN '[deleted]' ELSE display_name END) AS display_name,
    	(CASE WHEN deleted_at IS NOT NULL THEN '[deleted]' ELSE email END) AS email
    FROM users`,

	`CREATE VIEW IF NOT EXISTS user_sessions AS
		SELECT u.id, u.username, u.email, u.password, u.deleted_at, s.id AS session_id, s.created_at AS session_created, s.expiry, s.lat, s.lon, s.ip_address
	FROM users u
	JOIN sessions s ON u.id = s.user_id`,

	`CREATE VIEW IF NOT EXISTS vote_view AS
		SELECT post_id,
		SUM(CASE WHEN vote = 1 THEN 1 ELSE 0 END) as positive,
		SUM(CASE WHEN vote = -1 THEN 1 ELSE 0 END) as negative
	FROM votes
	GROUP BY post_id`,

	`CREATE VIEW IF NOT EXISTS children_view AS
		SELECT y.id, 
		(SELECT COUNT(p.id) FROM posts p WHERE p.parent LIKE CONCAT(y.parent, '/', y.id, '%')) AS num_children
	FROM posts y`,

	`CREATE VIEW IF NOT EXISTS post_view AS
		SELECT p.*, pu.username, pu.display_name, IFNULL(vv.positive, 0) as positive, IFNULL(vv.negative, 0) as negative, 
		u.id AS vote_user_id, u.username as vote_username, IFNULL(v.vote, 0) AS vote_state, c.num_children,
		IFNULL(((positive + 1.9208) / (positive + negative) - 1.96 * SQRT((positive * negative) / (positive + negative) + 0.9604) /(positive + negative)) / (1 + 3.8416 / (positive + negative)), 0) AS score
	FROM posts_live p 
	JOIN users_live pu ON p.user_id = pu.id
	LEFT JOIN vote_view vv ON vv.post_id = p.id
	CROSS JOIN users u
	LEFT JOIN votes v ON u.id = v.user_id AND p.id = v.post_id
	JOIN children_view c ON p.id = c.id
	WHERE u.deleted_at IS NULL`,
}
