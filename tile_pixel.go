/*
https://www.maptiler.com/google-maps-coordinates-tile-bounds-projection/
https://wiki.openstreetmap.org/wiki/Slippy_map_tilenames

*/

package transform

import (
	"math"
)

// ZXYtoWebMCBound SlippyMap->球面墨卡托边界
func ZXYtoWebMCBound(z, x, y float64) [4]int {
	east, north := ZXYtoWebMC(z, x, y)
	west, south := ZXYtoWebMC(z, x+1, y+1)
	return [4]int{east, south, west, north}
}

// ZXYtoWebMC SlippyMap->球面墨卡托
func ZXYtoWebMC(z, x, y float64) (xMc, yMc int) {
	res := mc * 2 / math.Pow(2, z)
	xMc = int(x*res - mc)
	yMc = int(-y*res + mc)
	return
}

// ZXYtoWGS84Bound SlippyMap->WGS84坐标系边界
func ZXYtoWGS84Bound(z, x, y int) [4]float64 {
	east, north := ZXYtoWGS84(z, x, y)
	west, south := ZXYtoWGS84(z, x+1, y+1)
	return [4]float64{east, south, west, north}
}

// WGS84toZXY WGS84坐标系->SlippyMap
func WGS84toZXY(lon, lat float64, z uint64) (x, y int) {
	x = int(math.Floor((lon + 180.0) / 360.0 * (math.Exp2(float64(z)))))
	y = int(math.Floor((1.0 - math.Log(math.Tan(lat*math.Pi/180.0)+1.0/math.Cos(lat*math.Pi/180.0))/math.Pi) / 2.0 * (math.Exp2(float64(z)))))
	return
}

// ZXYtoWGS84 SlippyMap->WGS84坐标系
func ZXYtoWGS84(z, x, y int) (lon, lat float64) {
	n := math.Pi - 2.0*math.Pi*float64(y)/math.Exp2(float64(z))
	lon = float64(x)/math.Exp2(float64(z))*360.0 - 180.0
	lat = 180.0 / math.Pi * math.Atan(0.5*(math.Exp(n)-math.Exp(-n)))
	return
}
