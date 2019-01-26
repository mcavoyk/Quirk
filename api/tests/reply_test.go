package tests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"testing"

	"github.com/mcavoyk/quirk/api/server"

	"github.com/stretchr/testify/assert"
	baloo "gopkg.in/h2non/baloo.v3"
)

const (
	url  = "http://localhost:5005"
	base = "/api/v1"
)

var (
	api = baloo.New(url)
	ids []string
)

func TestHealth(t *testing.T) {
	assert.Nil(t, api.Get(base+"/health").
		Expect(t).
		Status(200).
		Type("json").
		Done())
}

func TestPostReplies(t *testing.T) {
	createPostTree(t, 3, "")

	idList1 := make([]string, len(ids))
	copy(idList1, ids)
	ids = make([]string, 0)

	createPostTree(t, 5, "")

	idList2 := make([]string, len(ids))
	copy(idList2, ids)
	ids = make([]string, 0)

	token, err := auth(0, 0)
	if err != nil {
		t.Fatalf("Error authentication with api: %s", err.Error())
	}

	assert.Nil(t, api.Get(base+"/post/"+idList1[0]+"/posts").
		SetHeader("Authorization", "bearer "+token).
		Expect(t).
		Status(200).
		Type("json").
		AssertFunc(getListID).
		Done())

	assert.Equal(t, len(ids), len(idList1)-1)
	ids = make([]string, 0)

	assert.Nil(t, api.Get(base+"/post/"+idList2[0]+"/posts").
		SetHeader("Authorization", "bearer "+token).
		Expect(t).
		Status(200).
		Type("json").
		AssertFunc(getListID).
		Done())

	assert.Equal(t, len(ids), len(idList2)-1)
	ids = make([]string, 0)
}

func getListID(res *http.Response, req *http.Request) error {
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("Invalid server response body: %s", err.Error())
	}
	var out []map[string]interface{}
	err = json.Unmarshal(body, &out)
	if err != nil {
		return fmt.Errorf("Invalid server json body: %s", err.Error())
	}

	for _, e := range out {
		ids = append(ids, e["ID"].(string))
	}
	return nil
}

func getSingleID(res *http.Response, req *http.Request) error {
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("Invalid server response body: %s", err.Error())
	}
	var out map[string]interface{}
	err = json.Unmarshal(body, &out)
	if err != nil {
		return fmt.Errorf("Invalid server json body: %s", err.Error())
	}
	ids = append(ids, out["ID"].(string))
	return nil
}

// createPostTree creates a post of trees to try and replicate a real
// post thread. Height specifies how many posts to chain together, where 1 would result
// in a single post with no replies. Reply should be an empty string to start the tree
// or another postID to create a tree off of that post
func createPostTree(t *testing.T, height int, reply string) {
	if height == 0 {
		return
	}
	token, err := auth(0, 0)
	if err != nil {
		t.Fatalf("Error authentication with api: %s", err.Error())
	}
	if reply != "" {
		reply = "/" + reply + "/post"
	}
	post1 := server.Post{"Test content", "public", 0.0, 0.0}
	assert.Nil(t, api.Post(base+"/post"+reply).
		SetHeader("Authorization", "bearer "+token).
		JSON(post1).
		Expect(t).
		Status(200).
		Type("json").
		AssertFunc(getSingleID).
		Done())

	newReply := ids[len(ids)-1]
	children := rand.Intn(6)

	for i := 1; i <= children; i++ {
		createPostTree(t, height-1, newReply)
	}
}

func auth(lat, lon float64) (string, error) {
	resp, err := http.Get(url + base + fmt.Sprintf("/auth/token?lat=%f&lon=%f", lat, lon))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var jsonResp map[string]interface{}
	err = json.Unmarshal(body, &jsonResp)
	if err != nil {
		return "", err
	}
	tokenInterface, ok := jsonResp["token"]
	if !ok {
		return "", fmt.Errorf("error authorizing")
	}
	return tokenInterface.(string), nil
}
