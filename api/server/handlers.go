package server

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/mcavoyk/quirk/api/pkg/location"

	"github.com/gin-gonic/gin"
)

func (env *Env) HealthCheck(c *gin.Context) {
	ctx, cancel:= context.WithTimeout(context.Background(), time.Second)
	err := env.DB.PingContext(ctx)
	cancel()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "unhealthy",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "healthy"})
}

func noRoute(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"status": "Page not found",
	})
}

func (env *Env) selectQuery(c *gin.Context) {
	responseData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
		return
	}
	env.Log.Warnf("Executing read query: %s", string(responseData))
	rows, err := env.DB.Read.Query(string(responseData))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
		return
	}

	cols, err := rows.Columns()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
		return
	}

	// Result is your slice string.
	rawResult := make([][]byte, len(cols))
	result := make([]string, len(cols))

	dest := make([]interface{}, len(cols)) // A temporary interface{} slice
	for i, _ := range rawResult {
		dest[i] = &rawResult[i] // Put pointers to each string in the interface slice
	}

	allResults := [][]string{}
	for rows.Next() {
		err = rows.Scan(dest...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
			return
		}

		for i, raw := range rawResult {
			result[i] = cols[i] + ": "
			if raw == nil {
				result[i] += "NULL"
			} else {
				result[i] += string(raw)
			}
		}
		allResults = append(allResults, result)
	}

	c.JSON(http.StatusOK, allResults)
}

func extractCoords(c *gin.Context) (*location.Point, error) {
	latStr := c.Query("lat")
	lonStr := c.Query("lon")
	lat, latErr := strconv.ParseFloat(latStr, 64)
	lon, lonErr := strconv.ParseFloat(lonStr, 64)

	if latErr != nil || lonErr != nil {
		return nil, fmt.Errorf("Invalid or missing latitude and longitude")
	}
	return &location.Point{Lat: lat, Lon: lon}, nil
}

