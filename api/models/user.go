package models

import "github.com/sirupsen/logrus"

type UserService interface {
	ReadUser(string) (*User, error)
	ReadUserByName(string) (*User, error)
	AddUser(*User) (*User, error)
	UpdateUser(*User) (*User, error)
	DeleteUser(string) error
}

// User
type User struct {
	Default
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Password    string `json:"password,omitempty"`
	Email       string `json:"email"`
}

const (
	SelectUser       = "SELECT * FROM users WHERE id = ?"
	SelectUserByName = "SELECT * FROM users WHERE username = ?"
	DeleteUserSoft   = "UPDATE users SET deleted_at = NOW() WHERE id = ?"
)

var (
	InsertUser = InsertValues("INSERT INTO users (id, username, display_name, password, email)")
)

func (s *SqlDB) ReadUser(id string) (*User, error) {
	user := new(User)
	err := s.read.Get(user, SelectUser, id)
	return user, err
}

func (s *SqlDB) ReadUserByName(name string) (*User, error) {
	user := new(User)
	err := s.read.Get(user, SelectUserByName, name)
	return user, err
}

func (s *SqlDB) AddUser(user *User) (*User, error) {
	user.ID = NewGUID()
	_, err := s.write.NamedExec(InsertUser, user)
	if err != nil {
		logrus.Error("SQLDB: Error while adding user: %s", err.Error())
		return nil, err
	}
	return s.ReadUser(user.ID)
}

func (s *SqlDB) UpdateUser(*User) (*User, error) {
	return nil, nil
}

func (s *SqlDB) DeleteUser(id string) error {
	_, err := s.write.Exec(DeleteUserSoft, id)
	return err
}
