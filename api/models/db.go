package models

import (
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

	db.AutoMigrate(&Post{}, &Vote{}, &User{})

	db.Exec("CREATE TRIGGER score_insert " +
		"AFTER INSERT ON votes " +
		"FOR EACH ROW " +
		"UPDATE posts p SET p.score = p.score + NEW.state WHERE p.id = NEW.post_id")

	db.Exec("CREATE TRIGGER score_update " +
		"AFTER UPDATE ON votes " +
		"FOR EACH ROW " +
		"UPDATE posts p SET p.score = p.score - OLD.state + NEW.state WHERE p.id = NEW.post_id")
	return &DB{db}, nil
}

func NewGUID() string {
	return ksuid.New().String()
}
