package tests

import (
	"testing"

	"github.com/mcavoyk/quirk/api/pkg/gfyid"
	"github.com/mcavoyk/quirk/api/pkg/location"
	"github.com/mcavoyk/quirk/api/server"
	"github.com/stretchr/testify/assert"
)

func TestPostLocation(t *testing.T) {
	radius := 40.0    // km
	gridStep := 0.035 // lat/lon amount to increment by when making posts
	centerLat := location.ToRadians(30.283826)
	centerLon := location.ToRadians(-97.732547)

	bPoints := location.BoundingPoints(&location.Point{centerLat, centerLon}, radius)
	minPoint := location.Point{location.ToDegrees(bPoints[0].Lat), location.ToDegrees(bPoints[0].Lon)}
	maxPoint := location.Point{location.ToDegrees(bPoints[1].Lat), location.ToDegrees(bPoints[1].Lon)}

	ids = make([]string, 0)
	latPosts := int((maxPoint.Lat - minPoint.Lat) / gridStep)
	lonPosts := int((maxPoint.Lon - minPoint.Lon) / gridStep)
	for i := 0; i < latPosts; i++ {
		for j := 0; j < lonPosts; j++ {
			createPost(t, minPoint.Lat+(float64(i)*gridStep), minPoint.Lon+(float64(j)*gridStep))
		}
	}

	assert.Equal(t, latPosts*lonPosts, len(ids))
}

func createPost(t *testing.T, lat, lon float64) {
	token, err := auth(lat, lon)
	if err != nil {
		t.Fatalf("Error authentication with api: %s", err.Error())
	}

	post := server.Post{gfyid.RandomID(), "public", lat, lon}
	assert.Nil(t, api.Post(url+"/post").
		SetHeader("Authorization", "bearer "+token).
		JSON(post).
		Expect(t).
		Status(200).
		Type("json").
		AssertFunc(getSingleID).
		Done())
}
