package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/baloo.v3"
)

var api = baloo.New("http://localhost:5005/")
const base = "/api"

func TestHealth(t *testing.T) {
	assert.Nil(t, api.Get(base + "/health").
		Expect(t).
		Status(200).
		Type("json").
		Done())
}

func dropTables(t *testing.T) {
}

func insertPosts(users, postsPerUser int) {

}
