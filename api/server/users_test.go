package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/mcavoyk/quirk/api/models"

	"github.com/mcavoyk/quirk/api/mocks"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)


func performRequest(handler http.Handler, method, url string, body map[string]string) *httptest.ResponseRecorder {
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	return w
}

func TestCreateUserAlreadyExists(t *testing.T) {
	username := "testUser417"

	store := &mocks.Store{}
	store.On("Read", mock.AnythingOfType("*[]models.User"), models.SelectUserByName, username).
		Return(nil)

	router := NewRouter(store, viper.New())
	body := map[string]string{"username": username}
	w := performRequest(router, http.MethodPost, ApiV1+"/user", body)

	store.AssertNumberOfCalls(t, "Read", 1)
	store.AssertNumberOfCalls(t, "Write", 0)
	store.AssertNumberOfCalls(t, "Exec", 0)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateUserRandom(t *testing.T) {
	store := &mocks.Store{}

	randomUsername := ""
	newUser := models.User{}
	store.On("Read",  mock.AnythingOfType("*[]models.User"), models.SelectUserByName, mock.MatchedBy(func(username string) bool {
		randomUsername = username
		return true
	})).Return(fmt.Errorf("user not found")).Once()
	store.On("Read", mock.AnythingOfType("*[]models.User"), models.SelectUserByName, mock.MatchedBy(func(username string) bool {
		return randomUsername == username
	})).Return(nil).Run(func(args mock.Arguments) {
		*(args[0].(*[]models.User)) = []models.User{newUser}
	})
	store.On("Write", models.InsertValues(models.InsertUser), mock.MatchedBy(func(insert *models.User) bool {
		newUser = *insert
		return true
	})).Return(nil)


	router := NewRouter(store, viper.New())
	w := performRequest(router, http.MethodPost, ApiV1+"/user", nil)

	body, _ := ioutil.ReadAll(w.Body)
	var createdUser models.User
	_ = json.Unmarshal(body, &createdUser)
	store.AssertNumberOfCalls(t, "Read", 2)
	store.AssertNumberOfCalls(t, "Write", 1)
	store.AssertNumberOfCalls(t, "Exec", 0)
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.NotZero(t, createdUser.Password, "Should not return empty password")
	assert.NotEqual(t, newUser.Password, createdUser.Password, "Should not return hashed password")
	createdUser.Password = newUser.Password
	assert.True(t, reflect.DeepEqual(newUser, createdUser))
}

func TestCreateUser(t *testing.T) {
	body := map[string]string{"username": "testUsername417", "password": "hunter2"}

	store := &mocks.Store{}
	newUser := models.User{}
	store.On("Read",  mock.AnythingOfType("*[]models.User"), models.SelectUserByName, body["username"]).Return(fmt.Errorf("user not found")).Once()
	store.On("Read", mock.AnythingOfType("*[]models.User"), models.SelectUserByName, body["username"]).Return(nil).Run(func(args mock.Arguments) {
		*(args[0].(*[]models.User)) = []models.User{newUser}
	})
	store.On("Write", models.InsertValues(models.InsertUser), mock.MatchedBy(func(insert *models.User) bool {
		newUser = *insert
		return true
	})).Return(nil)


	router := NewRouter(store, viper.New())
	w := performRequest(router, http.MethodPost, ApiV1+"/user", body)

	respBody, _ := ioutil.ReadAll(w.Body)
	var createdUser models.User
	_ = json.Unmarshal(respBody, &createdUser)
	store.AssertNumberOfCalls(t, "Read", 2)
	store.AssertNumberOfCalls(t, "Write", 1)
	store.AssertNumberOfCalls(t, "Exec", 0)
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Zero(t, createdUser.Password, "Should return empty password field")
	createdUser.Password = newUser.Password
	assert.True(t, reflect.DeepEqual(newUser, createdUser))
}