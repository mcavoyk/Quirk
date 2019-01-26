package models

import (
	"fmt"
	"time"

	"github.com/mcavoyk/quirk/api/location"
)

type User struct {
	ID        string
	CreatedAt time.Time `gorm:"index:createdAt"`
	UsedAt    time.Time `gorm:"index:usedAt"`
	IP        string
	Lat       float64 `gorm:"index:latitude"`
	Lon       float64 `gorm:"index:longitude"`
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

func (db *DB) UsersByDistance(lat, lon float64) int {
	points := location.BoundingPoints(&location.Point{Lat: lat, Lon: lon}, Distance)
	minLat := points[0].Lat
	minLon := points[0].Lon
	maxLat := points[1].Lat
	maxLon := points[1].Lon

	row := db.Raw("SELECT COUNT(*) FROM users WHERE "+byDistance,
		minLat, maxLat, minLon, maxLon,
		lat, lat, lon, Distance/location.EarthRadius).Row()

	var count int
	err := row.Scan(&count)
	if err != nil {
		fmt.Printf("SQL Error: %s\n", err.Error())
		return 0
	}
	return count
}
