package models

import (
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/segmentio/ksuid"
)

type DB struct {
	*sqlx.DB
	log *logrus.Logger
}

// InitDB panics if unable to establish connection with DB
func InitDB(connection string) (*DB, error) {
	db, err := sqlx.Connect("mysql", connection + "?parseTime=True&charset=utf8mb4&collation=utf8mb4_unicode_ci")
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


// Database schema
var schema = []string{
	`CREATE DATABASE IF NOT EXISTS quirk CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci;`,

	`USE quirk;`,
	`CREATE TABLE IF NOT EXISTS metadata (
  		name VARCHAR(255),
		version VARCHAR(255),
  		updatedAt TIMESTAMP NOT NULL DEFAULT NOW() ON UPDATE NOW(),
  		PRIMARY KEY (name)
	);`,

	`CREATE TABLE IF NOT EXISTS users (
 		id VARCHAR(255) PRIMARY KEY NOT NULL,
  		createdAt TIMESTAMP NOT NULL DEFAULT NOW(),
		updatedAt TIMESTAMP NOT NULL DEFAULT NOW() ON UPDATE NOW(),
  		deletedAt TIMESTAMP NULL,
  		name VARCHAR(255) NOT NULL,
  		ip VARCHAR(255),
  		lat DOUBLE NOT NULL,
  		lon DOUBLE NOT NULL,

 		INDEX idxLatLon (lat, lon),
		INDEX idxName (name),
		INDEX idxCreated (createdAt),
		INDEX idxUpdated (updatedAt),
		INDEX idxDeleted (deletedAt)
	);`,

	`CREATE TABLE IF NOT EXISTS posts (
  		id VARCHAR(255) PRIMARY KEY NOT NULL,
		createdAt TIMESTAMP NOT NULL DEFAULT NOW(),
		updatedAt TIMESTAMP NOT NULL DEFAULT NOW() ON UPDATE NOW(),
		deletedAt TIMESTAMP NULL,
		lat DOUBLE NOT NULL,
		lon DOUBLE NOT NULL,
		userID VARCHAR(255) NOT NULL,
		parentID VARCHAR(1023) NOT NULL DEFAULT '',
		accessType ENUM('public', 'private') NOT NULL DEFAULT 'public',
		content TEXT NOT NULL,

		INDEX idxLatLon (lat, lon),
		INDEX idxCreated (createdAt),
		INDEX idxDeleted (deletedAt),

		FOREIGN KEY fk_users (userID)
			REFERENCES users(id)
			ON DELETE CASCADE
	);`,

	`CREATE TABLE IF NOT EXISTS votes (
  		postID VARCHAR(255) NOT NULL,
  		userID VARCHAR(255) NOT NULL,
  		vote TINYINT NOT NULL,
  		updatedAt TIMESTAMP NOT NULL DEFAULT NOW() ON UPDATE NOW(),
  		PRIMARY KEY (postID, userID),

  		INDEX idxPost (postID),
  		INDEX idxUser (userID),
		INDEX idxVote (vote),

  		FOREIGN KEY fkUser (userID)
    		REFERENCES users(id)
    		ON DELETE CASCADE,

  		FOREIGN KEY fkPost (postID)
  		  REFERENCES posts(id)
    		ON DELETE CASCADE
	);`,

	`CREATE VIEW IF NOT EXISTS voteView AS
		SELECT postID,
		SUM(CASE WHEN vote = 1 THEN 1 ELSE 0 END) as positive,
		SUM(CASE WHEN vote = -1 THEN 1 ELSE 0 END) as negative
	FROM votes
	GROUP BY postID`,

	`CREATE VIEW IF NOT EXISTS postView AS
		SELECT p.id, p.createdAt, p.updatedAt, p.lat, p.lon, p.userID, p.accessType, p.parentID, 
		p.content, u.name, v.positive, v.negative
	FROM posts p 
	JOIN users u ON p.userID = u.id
	JOIN voteView v ON v.postID = p.id`,
}

