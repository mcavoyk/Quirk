package models

import (
	"time"
)

type User struct {
	ID        string
	CreatedAt time.Time `gorm:"index:createdAt"`
	UsedAt    time.Time `gorm:"index:usedAt"`
	IP        string
}

func (db *DB) UserInsert(user *User) string {
	user.ID = NewGUID()
	user.CreatedAt = time.Now()
	user.UsedAt = user.CreatedAt
	db.Create(user)
	return user.ID
}

func (db *DB) UserGet(id string) *User {
	user := new(User)
	db.Where("ID = ?", id).First(user)
	return user
}

func (db *DB) UserUpdate(user *User) {
	if user.ID == "" {
		return
	}
	db.Model(user).Updates(user)
	return
}
