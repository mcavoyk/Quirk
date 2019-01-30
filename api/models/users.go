package models

import (
	"fmt"
	"time"

	"github.com/mcavoyk/quirk/api/location"
)

type User struct {
	ID        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	IP        string
	Lat       float64
	Lon       float64
}

const insertUser = "INSERT INTO users (id, name, ip, lat, lon) VALUES (?, ?, ?, ?, ?)"

func (db *DB) UserInsert(user *User) string {
	user.ID = NewGUID()
	_, err := db.NamedExec(insertUser, user)
	db.log.Errorf("User insert failed: %s", err.Error())
	return user.ID
}

func (db *DB) UserGet(id string) *User {
	user := new(User)
	_ = db.QueryRow("")
	return user
}

func (db *DB) UserUpdate(user *User) {
	if user.ID == "" {
		return
	}
	return
}

func (db *DB) UsersByDistance(lat, lon float64) int {
	points := location.BoundingPoints(&location.Point{Lat: lat, Lon: lon}, Distance)
	minLat := points[0].Lat
	minLon := points[0].Lon
	maxLat := points[1].Lat
	maxLon := points[1].Lon

	row := db.QueryRowx("SELECT COUNT(*) FROM users WHERE "+byDistance,
		minLat, maxLat, minLon, maxLon,
		lat, lat, lon, Distance/location.EarthRadius)

	var count int
	err := row.Scan(&count)
	if err != nil {
		fmt.Printf("SQL Error: %s\n", err.Error())
		return 0
	}
	return count
}
