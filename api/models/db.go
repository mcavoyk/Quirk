package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type DB struct {
	*gorm.DB
}

func InitDB(connection string) (*DB, error) {
	db, err := gorm.Open("mysql", connection)
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&Post{}, &Comment{})
	return &DB{db}, nil
}
