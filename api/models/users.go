package models

import (
	"fmt"

	"github.com/mcavoyk/quirk/api/location"
)

type User struct {
	Default
	Name     string  `json:"name"`
	Password string  `json:"-"`
	Email    string  `json:"email"`
	IP       string  `json:"ip_address"`
	Lat      float64 `json:"lat"`
	Lon      float64 `json:"lon"`
}

const insertUser = "INSERT INTO users (id, name, ip_address, lat, lon) VALUES (:id, :name, :ip_address, :lat, :lon)"

func (db *DB) InsertUser(user *User) *User {
	user.ID = NewGUID()
	_, err := db.NamedExec(insertUser, user)
	if err != nil {
		db.log.Errorf("Insert user failed: %s", err.Error())
	}
	return db.GetUser(user.ID)
}

func (db *DB) GetUser(id string) *User {
	user := new(User)
	err := db.Get(user, "SELECT * FROM users WHERE id=?", id)
	if err != nil {
		db.log.Errorf("Get user failed: %s", err.Error())
	}
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
