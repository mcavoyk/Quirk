package models

import (
	"github.com/jmoiron/sqlx"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/segmentio/ksuid"
)

type DB struct {
	*sqlx.DB
}

// InitDB panics if unable to establish connection with DB
func InitDB(connection string) (*DB, error) {
	db, err := sqlx.Open("mysql", connection + "?parseTime=True&charset=utf8mb4&collation=utf8mb4_unicode_ci")
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Minute * 10)

	for _, stmt := range schema {
		_, err = db.Exec(stmt)
		if err != nil {
			return nil, err
		}
	}
	return &DB{db}, nil
}

func NewGUID() string {
	return ksuid.New().String()
}


// Database schema
var schema = []string{
	`CREATE DATABASE IF NOT EXISTS quirk_db CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci;`,

	`USE quirk_db;`,
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
  		name VARCHAR(255) NOT NULL,
  		ip VARCHAR(255),
  		lat DOUBLE NOT NULL,
  		lon DOUBLE NOT NULL,

 		INDEX idx_latlon (lat, lon)
	);`,

	`CREATE TABLE IF NOT EXISTS posts (
  		id VARCHAR(255) PRIMARY KEY NOT NULL,
 		 created_at TIMESTAMP NOT NULL DEFAULT NOW(),
 		 updated_at TIMESTAMP NOT NULL DEFAULT NOW() ON UPDATE NOW(),
 		 deleted_at TIMESTAMP NULL,
 		 lat DOUBLE NOT NULL,
 		 lon DOUBLE NOT NULL,
 		 user_id VARCHAR(255) NOT NULL,
 		 parent_id VARCHAR(1023) NOT NULL,
 		 access_type ENUM('public', 'private') NOT NULL,
 		 content TEXT NOT NULL,

 		 INDEX idx_latlon (lat, lon)
	);`,

	`CREATE TABLE IF NOT EXISTS votes (
  		post_id VARCHAR(255) NOT NULL,
  		user_id VARCHAR(255) NOT NULL,
  		vote TINYINT NOT NULL,
  		updated_at TIMESTAMP NOT NULL DEFAULT NOW() ON UPDATE NOW(),
  		PRIMARY KEY (post_id, user_id),

  		INDEX idx_post (post_id),
  		INDEX idx_user (user_id),

  		FOREIGN KEY fk_user (user_id)
    		REFERENCES users(id)
    		ON DELETE CASCADE,

  		FOREIGN KEY fk_post (post_id)
  		  REFERENCES posts(id)
    		ON DELETE CASCADE
	);`,
}

