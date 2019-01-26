package models

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/segmentio/ksuid"
)

type DB struct {
	*gorm.DB
}

func InitDB(connection string) (*DB, error) {
	db, err := gorm.Open("mysql", connection)
	if err != nil {
		return nil, err
	}

	db.DB().SetMaxOpenConns(10)
	db.DB().SetMaxIdleConns(5)
	db.DB().SetConnMaxLifetime(time.Minute * 10)

	db.AutoMigrate(&Post{}, &Vote{}, &User{})

	db.Exec("CREATE TRIGGER score_insert " +
		"AFTER INSERT ON votes " +
		"FOR EACH ROW BEGIN " +
		"IF NEW.state = 1 THEN UPDATE posts p SET p.positive = p.positive + 1 WHERE p.id = NEW.post_id; " +
		"ELSEIF NEW.state = -1 THEN UPDATE posts p SET p.negative = p.negative + 1 WHERE p.id = NEW.post_id; END IF; END")

	db.Exec("CREATE TRIGGER score_update " +
		"AFTER UPDATE ON votes " +
		"FOR EACH ROW BEGIN " +
		"IF NEW.state = 1 THEN UPDATE posts p SET p.positive = p.positive + 1, p.negative = p.negative + OLD.state WHERE p.id = NEW.post_id; " +
		"ELSEIF NEW.state = -1 THEN UPDATE posts p SET p.negative = p.negative + 1, p.positive = p.positive - old.state WHERE p.id = NEW.post_id; " +
		"ELSEIF NEW.state = 0 AND OLD.state = 1 THEN UPDATE posts p SET p.positive = p.positive - 1 WHERE p.id = NEW.post_id; " +
		"ELSEIF NEW.state = 0 AND OLD.state = -1 THEN UPDATE posts p SET p.negative = p.negative - 1 WHERE p.id = NEW.post_id; END IF; END")
	return &DB{db}, nil
}

func NewGUID() string {
	return ksuid.New().String()
}
