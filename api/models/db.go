package models

import (
	"database/sql/driver"
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

// InitDB panics if unable to establish connection with DB
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
		user_id VARCHAR(255) NOT NULL,
		parent_id VARCHAR(1023) NOT NULL DEFAULT '',
		access_type ENUM('public', 'private') NOT NULL DEFAULT 'public',
		content TEXT NOT NULL,

		INDEX idx_parent (parent_id),
		INDEX idx_latlon (lat, lon),
		INDEX idx_created (created_at),
		INDEX idx_deleted (updated_at),

		FOREIGN KEY fk_user_posts (user_id)
			REFERENCES users(id)
			ON DELETE CASCADE
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
    		ON DELETE CASCADE,

  		FOREIGN KEY fk_post_votes (post_id)
  		  REFERENCES posts(id)
    		ON DELETE CASCADE
	);`,

	`CREATE VIEW IF NOT EXISTS user_view AS
		SELECT u.id, u.username, u.email, u.password, u.deleted_at, s.id AS session_id, s.created_at AS session_created, s.expiry, s.lat, s.lon, s.ip_address
	FROM users u
	JOIN sessions s ON u.id = s.user_id`,

	`CREATE VIEW IF NOT EXISTS vote_view AS
		SELECT post_id,
		SUM(CASE WHEN vote = 1 THEN 1 ELSE 0 END) as positive,
		SUM(CASE WHEN vote = -1 THEN 1 ELSE 0 END) as negative
	FROM votes
	GROUP BY post_id`,

	`CREATE VIEW IF NOT EXISTS post_view AS
		SELECT p.id, p.created_at, p.updated_at, p.lat, p.lon, p.user_id, p.access_type, p.parent_id, 
		p.content, u.username, u.display_name, v.positive, v.negative
	FROM posts p 
	JOIN users u ON p.user_id = u.id
	JOIN vote_view v ON v.post_id = p.id`,
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
