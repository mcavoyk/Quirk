package server

import (
	"github.com/gin-gonic/gin"
	"github.com/mcavoyk/quirk/api/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

func init() {
	router := gin.Default()
	loadRoutes(router, &Env{DB: models.Store{}})
}

func TestCreateUser(t *testing.T) {
	req, _ := http.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Env{}.CreateUser)
}