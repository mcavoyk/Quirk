package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/mcavoyk/quirk/api/location"

	"github.com/gin-gonic/gin"
)

func (env *Env) HealthCheck(c *gin.Context) {
	err := env.DB.DB.DB().Ping()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "unhealthy",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "healthy"})
}

func noRoute(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"message": "Page not found",
	})
}

func extractCoords(c *gin.Context) (*location.Point, error) {
	latStr := c.Query("lat")
	lonStr := c.Query("lon")
	lat, latErr := strconv.ParseFloat(latStr, 64)
	lon, lonErr := strconv.ParseFloat(lonStr, 64)

	if latErr != nil || lonErr != nil {
		return nil, fmt.Errorf("Invalid or missing latitude and longitude")
	}
	return &location.Point{Lat: location.ToRadians(lat), Lon: location.ToRadians(lon)}, nil
}
