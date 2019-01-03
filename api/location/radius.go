package location

import "math"

// Point on Earth measured in Degrees
type Point struct {
	Lat float64
	Lon float64
}

const (
	EarthRadius = 6371.01
	MinLat      = (-math.Pi) / 2 // -PI/2
	MaxLat      = math.Pi / 2    //  PI/2
	MinLon      = -math.Pi       // -PI
	MaxLon      = math.Pi        //  PI
)

// BoundingPoints returns the bounding coordinates of all points on the surface
// a sphere that have a great circle distance to the point represented
// by this GeoLocation instance that is less or equal to the distance argument
// Created based on http://janmatuschek.de/LatitudeLongitudeBoundingCoordinate
func BoundingPoints(point *Point, distance float64) []Point {
	radDist := distance / EarthRadius
	radLat := point.Lat
	radLon := point.Lon

	minLat := radLat - radDist
	maxLat := radLat + radDist

	var minLon, maxLon float64
	if minLat > MinLat && maxLat < MaxLat {
		deltaLon := math.Asin(math.Sin(radDist)) / math.Cos(radLat)
		minLon = radLon - deltaLon
		if minLon < MinLon {
			minLon += 2 * math.Pi
		}
		maxLon = radLon + deltaLon
		if maxLon > MaxLon {
			maxLon -= 2 * math.Pi
		}
	} else {
		// a pole is within the distance
		minLat = math.Max(minLat, MinLat)
		maxLat = math.Min(maxLat, MaxLat)
		minLon = MinLon
		maxLon = MaxLon
	}

	boundPoints := make([]Point, 2)
	boundPoints[0] = Point{minLat, minLon}
	boundPoints[1] = Point{maxLat, maxLon}
	return boundPoints
}

func ToRadians(x float64) float64 {
	return x * (math.Pi / 180)
}

func ToDegrees(x float64) float64 {
	return x * (180 / math.Pi)
}
