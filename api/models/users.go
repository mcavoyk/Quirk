package models

import (
	"time"

	"github.com/sirupsen/logrus"

	"github.com/mcavoyk/quirk/api/pkg/location"
)

type User struct {
	Default
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Password    string `json:"password,omitempty"`
	Email       string `json:"email"`
}

type Session struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Expiry    time.Time `json:"expiry"`
	IP        string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	Lat       float64   `json:"lat"`
	Lon       float64   `json:"lon"`
}

const insertUser = "INSERT INTO users (id, username, display_name, password, email)"
const insertSession = "INSERT INTO sessions (id, user_id, expiry, ip_address, user_agent, lat, lon)"

func (db *DB) InsertUser(user *User) (*User, error) {
	user.ID = NewGUID()
	logrus.Infof("Creating user: %s", user.ID)
	sqlStmt := InsertValues(insertUser)
	_, err := db.NamedExec(sqlStmt, user)
	if err != nil {
		logrus.Debugf("Insert user SQL: %s", sqlStmt)
		logrus.Warnf("Insert user failed: %s", err.Error())
		return nil, err
	}
	return db.GetUser(user.ID)
}

func (db *DB) GetUser(id string) (*User, error) {
	user := new(User)
	err := db.Unsafe().Get(user, "SELECT * FROM users WHERE id=?", id)
	if err != nil {
		logrus.Debugf("Get user failed: %s", err.Error())
		return nil, err
	}
	return user, nil
}

func (db *DB) GetUserByName(username string) (*User, error) {
	user := new(User)
	err := db.Get(user, "SELECT * FROM users WHERE username=?", username)
	if err != nil {
		logrus.Debugf("Get user failed: %s", err.Error())
		return nil, err
	}

	return user, nil
}

func (db *DB) InsertSession(session *Session) (*Session, error) {
	session.ID = NewGUID()
	session.Lat, session.Lon = location.ToRadians(session.Lat), location.ToRadians(session.Lon)
	sqlStmt := InsertValues(insertSession)
	_, err := db.NamedExec(sqlStmt, session)
	if err != nil {
		logrus.Debugf("Insert session SQL: %s", sqlStmt)
		logrus.Warnf("Insert session failed: %s", err.Error())
		return nil, err
	}
	return db.GetSession(session.ID)
}

func (db *DB) GetSession(id string) (*Session, error) {
	session := new(Session)
	err := db.Get(session, "SELECT * FROM sessions WHERE id=?", id)
	if err != nil {
		logrus.Debugf("Get session failed: %s", err.Error())
		return nil, err
	}
	return session, nil
}

func (db *DB) GetUserBySession(sessionID string) (*User, error) {
	user := new(User)
	err := db.Unsafe().Get(user, "SELECT * FROM user_sessions WHERE session_id=?", sessionID)
	if err != nil {
		logrus.Errorf("Get user by session failed: %s", err.Error())
		return nil, err
	}

	return user, nil
}

func (db *DB) SessionUpdate(session *Session) {
	session.Lat, session.Lon = location.ToRadians(session.Lat), location.ToRadians(session.Lon)
	_, err := db.NamedExec("UPDATE sessions SET ip_address = :ip_address, lat = :lat, lon = :lon WHERE id = :id", session)
	if err != nil {
		logrus.Errorf("Update session failed: %s", err.Error())
	}
}

func (db *DB) UpdateUser(user *User) (*User, error) {
	sqlStmt := "UPDATE users " + createSet(*user) + " WHERE id = ?"
	logrus.Debugf("Update post SQL: %s", sqlStmt)
	_, err := db.Exec(sqlStmt, user.ID)
	if err != nil {
		logrus.Debugf("Update user failed: %s", err.Error())
		//logrus.Debugf("Update post SQL: %s", sqlStmt)
		return nil, err
	}
	return db.GetUser(user.ID)
}

func (db *DB) DeleteUser(id string) error {
	var returnedErr error
	_, err := db.Exec("UPDATE users SET deleted_at = NOW() WHERE id = ?", id)
	if err != nil {
		logrus.Errorf("Delete user failed: %s", err.Error())
		returnedErr = err
	}
	_, err = db.Exec("DELETE FROM sessions WHERE user_id = ?", id)
	if err != nil {
		logrus.Errorf("Delete sessions failed: %s", err.Error())
		returnedErr = err
	}
	return returnedErr
}

/*
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
*/
