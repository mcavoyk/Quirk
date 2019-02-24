package server

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/mcavoyk/quirk/api/pkg/location"

	"github.com/mcavoyk/quirk/api/models"
	"github.com/stretchr/testify/mock"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestCreatePostInvalidContent(t *testing.T) {
	store := authStore()

	router := NewRouter(store, viper.New())
	body := map[string]interface{}{"content": "dank memes", "access_type": "public", "lat": 0.0, "lon": 0.0}
	w := performRequest(router, http.MethodPost, ApiV1+"/post", body)

	assertStore(t, store, authCalls(storeFunc{}))
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreatePostInvalidParent(t *testing.T) {
	store := authStore()
	parentID := "703"

	store.On("ReadOne", mock.AnythingOfType("*models.PostInfo"), models.SelectPostByUser, parentID, testUser).Return(fmt.Errorf("no parent post found"))

	router := NewRouter(store, viper.New())
	body := map[string]interface{}{"content": map[string]interface{}{"title": "dank memes"}, "access_type": "public", "lat": 0.0, "lon": 0.0}
	w := performRequest(router, http.MethodPost, ApiV1+"/post/"+parentID+"/post", body)

	assertStore(t, store, authCalls(storeFunc{ReadOne: 1}))
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreatePost(t *testing.T) {
	body := map[string]interface{}{"content": map[string]interface{}{"title": "dank memes"}, "access_type": "public", "lat": 5.0, "lon": 10.0}
	store := authStore()

	store.On("Write", models.InsertPost, mock.MatchedBy(func(insert *models.Post) bool {
		result := insert.UserID == testUser && insert.ID != "" && insert.Parent == "" && insert.AccessType == body["access_type"].(string)
		result = result && insert.Lat == location.ToRadians(body["lat"].(float64)) && insert.Lon == location.ToRadians(body["lon"].(float64))
		return result
	})).Return(nil)
	store.On("ReadOne", mock.AnythingOfType("*models.PostInfo"), models.SelectPostByUser, mock.AnythingOfType("string"), testUser).Return(nil)

	router := NewRouter(store, viper.New())
	w := performRequest(router, http.MethodPost, ApiV1+"/post", body)

	assertStore(t, store, authCalls(storeFunc{Write: 1, ReadOne: 1}))
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCreatePostReply(t *testing.T) {
	parentID := "1004"
	body := map[string]interface{}{"content": map[string]interface{}{"title": "dank memes"}, "access_type": "public", "lat": 5.0, "lon": 10.0}
	store := authStore()

	store.On("ReadOne", mock.AnythingOfType("*models.PostInfo"), models.SelectPostByUser, parentID, testUser).Return(nil).Once()
	store.On("Write", models.InsertPost, mock.MatchedBy(func(insert *models.Post) bool {
		result := insert.UserID == testUser && insert.ID != "" && insert.Parent == parentID && insert.AccessType == body["access_type"].(string)
		result = result && insert.Lat == location.ToRadians(body["lat"].(float64)) && insert.Lon == location.ToRadians(body["lon"].(float64))
		return result
	})).Return(nil)
	store.On("ReadOne", mock.AnythingOfType("*models.PostInfo"), models.SelectPostByUser, mock.AnythingOfType("string"), testUser).Return(nil)

	router := NewRouter(store, viper.New())
	w := performRequest(router, http.MethodPost, ApiV1+"/post/"+parentID+"/post", body)

	assertStore(t, store, authCalls(storeFunc{Write: 1, ReadOne: 2}))
	assert.Equal(t, http.StatusCreated, w.Code)
}

///		 Get post by ID		 \\\
func TestGetPostNotExist(t *testing.T) {
	store := authStore()
	postID := "703"

	store.On("ReadOne", mock.AnythingOfType("*models.PostInfo"), models.SelectPostByUser, postID, testUser).Return(fmt.Errorf("no post found"))

	router := NewRouter(store, viper.New())
	w := performRequest(router, http.MethodGet, ApiV1+"/post/"+postID, nil)

	assertStore(t, store, authCalls(storeFunc{ReadOne: 1}))
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetPost(t *testing.T) {
	store := authStore()
	postID := "703"

	store.On("ReadOne", mock.AnythingOfType("*models.PostInfo"), models.SelectPostByUser, postID, testUser).Return(nil)

	router := NewRouter(store, viper.New())
	w := performRequest(router, http.MethodGet, ApiV1+"/post/"+postID, nil)

	assertStore(t, store, authCalls(storeFunc{ReadOne: 1}))
	assert.Equal(t, http.StatusOK, w.Code)
}

///		 Get post by lat/lon	 \\\
func TestGetPostsByArea(t *testing.T) {

}

///		 Get post by parent		 \\\
func TestGetPostsParent(t *testing.T) {

}
