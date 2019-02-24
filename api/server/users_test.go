package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/mcavoyk/quirk/api/mocks"
	"github.com/mcavoyk/quirk/api/models"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const testToken = "sessionToken"
const testUser = "testUser417"

func performRequest(handler http.Handler, method, url string, body map[string]interface{}) *httptest.ResponseRecorder {
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", testToken))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	return w
}

func TestCreateUserAlreadyExists(t *testing.T) {
	username := "testUser417"

	store := &mocks.Store{}
	store.On("ReadOne", mock.AnythingOfType("*models.User"), models.SelectUserByName, username).
		Return(nil)

	router := NewRouter(store, viper.New())
	body := map[string]interface{}{"username": username}
	w := performRequest(router, http.MethodPost, ApiV1+"/user", body)

	assertStore(t, store, storeFunc{ReadOne: 1})
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateUserRandom(t *testing.T) {
	store := &mocks.Store{}

	randomUsername := ""
	newUser := models.User{}
	store.On("ReadOne", mock.AnythingOfType("*models.User"), models.SelectUserByName, mock.MatchedBy(func(username string) bool {
		randomUsername = username
		return true
	})).Return(fmt.Errorf("user not found")).Once()
	store.On("ReadOne", mock.AnythingOfType("*models.User"), models.SelectUserByName, mock.MatchedBy(func(username string) bool {
		return randomUsername == username
	})).Return(nil).Run(func(args mock.Arguments) {
		*(args[0].(*models.User)) = newUser
	})
	store.On("Write", models.InsertUser, mock.MatchedBy(func(insert *models.User) bool {
		newUser = *insert
		return insert.ID != ""
	})).Return(nil)

	router := NewRouter(store, viper.New())
	w := performRequest(router, http.MethodPost, ApiV1+"/user", nil)

	body, _ := ioutil.ReadAll(w.Body)
	var createdUser models.User
	_ = json.Unmarshal(body, &createdUser)

	assertStore(t, store, storeFunc{ReadOne: 2, Write: 1})
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.NotZero(t, createdUser.Password, "Should not return empty password")
	assert.NotEqual(t, newUser.Password, createdUser.Password, "Should not return hashed password")
	createdUser.Password = newUser.Password
	assert.True(t, reflect.DeepEqual(newUser, createdUser))
}

func TestCreateUser(t *testing.T) {
	body := map[string]interface{}{"username": "testUsername417", "password": "hunter2"}

	store := &mocks.Store{}
	newUser := models.User{}
	store.On("ReadOne", mock.AnythingOfType("*models.User"), models.SelectUserByName, body["username"]).Return(fmt.Errorf("user not found")).Once()
	store.On("ReadOne", mock.AnythingOfType("*models.User"), models.SelectUserByName, body["username"]).Return(nil).Run(func(args mock.Arguments) {
		*(args[0].(*models.User)) = newUser
	})
	store.On("Write", models.InsertUser, mock.MatchedBy(func(insert *models.User) bool {
		newUser = *insert
		return insert.ID != ""
	})).Return(nil)

	router := NewRouter(store, viper.New())
	w := performRequest(router, http.MethodPost, ApiV1+"/user", body)

	respBody, _ := ioutil.ReadAll(w.Body)
	var createdUser models.User
	_ = json.Unmarshal(respBody, &createdUser)

	assertStore(t, store, storeFunc{ReadOne: 2, Write: 1})
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Zero(t, createdUser.Password, "Should return empty password field")
	createdUser.Password = newUser.Password
	assert.True(t, reflect.DeepEqual(newUser, createdUser))
}

func TestLoginNoUser(t *testing.T) {
	body := map[string]interface{}{"username": "testUsername417", "password": "hunter2", "lat": 42, "lon": 42}

	store := &mocks.Store{}
	store.On("ReadOne", mock.AnythingOfType("*models.User"), models.SelectUserByName, body["username"]).Return(fmt.Errorf("user not found"))

	router := NewRouter(store, viper.New())
	w := performRequest(router, http.MethodPost, ApiV1+"/user/login", body)

	assertStore(t, store, storeFunc{ReadOne: 1})
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLoginUnauthorized(t *testing.T) {
	body := map[string]interface{}{"username": "testUsername417", "password": "hunter2", "lat": 42.0, "lon": 42.0}
	UserID := "666"

	store := &mocks.Store{}
	store.On("ReadOne", mock.AnythingOfType("*models.User"), models.SelectUserByName, body["username"]).Return(nil).Run(func(args mock.Arguments) {
		*(args[0].(*models.User)) = models.User{Default: models.Default{ID: UserID}, Password: "notHunter2Hash"}
	})

	router := NewRouter(store, viper.New())
	w := performRequest(router, http.MethodPost, ApiV1+"/user/login", body)

	assertStore(t, store, storeFunc{ReadOne: 1})
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLoginAuthorized(t *testing.T) {
	body := map[string]interface{}{"username": "testUsername417", "password": "hunter2", "lat": 42.0, "lon": 42.0}
	UserID := "666"

	store := &mocks.Store{}
	store.On("ReadOne", mock.AnythingOfType("*models.User"), models.SelectUserByName, body["username"]).Return(nil).Run(func(args mock.Arguments) {
		*(args[0].(*models.User)) = models.User{Default: models.Default{ID: UserID}, Password: "$2a$10$9IDtcGMy7SujzKJCoMibwO/5BXvYERn9GdVwtWXiV0ow2p1RPyakW"}
	})
	store.On("Write", models.InsertSession, mock.MatchedBy(func(insert *models.Session) bool {
		return insert.ID != "" && insert.UserID == UserID && insert.Lat == body["lat"].(float64) && insert.Lon == body["lon"].(float64)
	})).Return(nil)

	router := NewRouter(store, viper.New())
	w := performRequest(router, http.MethodPost, ApiV1+"/user/login", body)

	respBody, _ := ioutil.ReadAll(w.Body)
	var creds map[string]string
	_ = json.Unmarshal(respBody, &creds)
	expiry, _ := time.Parse(time.RFC3339, creds["expiry"])

	assertStore(t, store, storeFunc{ReadOne: 1, Write: 1})
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 2, len(creds))
	assert.NotZero(t, creds["token"], "Login session token should exist")
	assert.WithinDuration(t, time.Now().Add(tokenExpiry), expiry, time.Second, "Token expiry should be set")
}

// creates a mock with valid auth middleware configured
func authStore() *mocks.Store {
	store := &mocks.Store{}
	store.On("ReadOne", mock.AnythingOfType("*models.Session"), models.SelectSession, testToken).
		Return(nil).Run(func(args mock.Arguments) {
		*(args[0].(*models.Session)) = models.Session{UserID: testUser}
	}).Once()
	store.On("Write", models.UpdateSession, mock.AnythingOfType("models.Session")).
		Return(nil).Once()

	return store
}

func authCalls(store storeFunc) storeFunc {
	store.ReadOne = store.ReadOne + 1
	store.Write = store.Write + 1
	return store
}

func assertStore(t *testing.T, store *mocks.Store, funcCalls storeFunc) {
	var storeFuncs map[string]int
	out, _ := json.Marshal(funcCalls)
	_ = json.Unmarshal(out, &storeFuncs)
	for k, v := range storeFuncs {
		store.AssertNumberOfCalls(t, k, v)
	}
}

type storeFunc struct {
	Read    int
	ReadOne int
	Write   int
	Exec    int
}
