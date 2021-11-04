package transform

import (
	"math"
)

// Distance 经纬度距离计算，返回米
func Distance(lon1, lat1, lon2, lat2 float64) float64 {
	lat1r := lat1 * math.Pi / 180
	lat2r := lat2 * math.Pi / 180
	latRd := (lat2 - lat1) * math.Pi / 180
	lonRd := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(latRd/2)*math.Sin(latRd/2) + math.Cos(lat1r)*math.Cos(lat2r)*math.Sin(lonRd/2)*math.Sin(lonRd/2)
	return 2 * 6371e3 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}
