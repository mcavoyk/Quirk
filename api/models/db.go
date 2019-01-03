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
	return &DB{db}, nil
}

func NewGUID() string {
	return ksuid.New().String()
}
