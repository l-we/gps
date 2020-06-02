/*
https://en.wikipedia.org/wiki/Restrictions_on_geographic_data_in_China
https://github.com/sshuair/coord-convert
https://github.com/googollee/eviltransform/blob/master/go/transform.go
WGS84坐标系：即地球坐标系，国际上通用的坐标系。
GCJ02坐标系：即火星坐标系，WGS84坐标系经加密后的坐标系。
BD09坐标系：即百度坐标系，GCJ02坐标系经加密后的坐标系。
*/

package transform

import (
	"math"
)

const (
	xPi = math.Pi * 3000.0 / 180.0
	a   = 6378245.0
	ee  = 0.00669342162296594323
	mc  = 20037508.34

	threshold = 0.000001
)

func inChina(lon, lat float64) bool {
	return !(lon > 72.004 && lon < 137.8347 && lat > 0.8293 && lat < 55.8271)
}

// BD09toGCJ02 百度坐标系->火星坐标系
func BD09toGCJ02(lon, lat float64) (float64, float64) {
	x := lon - 0.0065
	y := lat - 0.006
	z := math.Sqrt(x*x+y*y) - 0.00002*math.Sin(y*xPi)
	theta := math.Atan2(y, x) - 0.000003*math.Cos(x*xPi)
	return z * math.Cos(theta), z * math.Sin(theta)
}

// GCJ02toBD09 火星坐标系->百度坐标系
func GCJ02toBD09(lon, lat float64) (float64, float64) {
	z := math.Sqrt(lon*lon+lat*lat) + 0.00002*math.Sin(lat*xPi)
	theta := math.Atan2(lat, lon) + 0.000003*math.Cos(lon*xPi)
	return z*math.Cos(theta) + 0.0065, z*math.Sin(theta) + 0.006
}

// WGS84toGCJ02 WGS84坐标系->火星坐标系
func WGS84toGCJ02(lon, lat float64) (float64, float64) {
	if inChina(lon, lat) {
		return lon, lat
	}
	dLat, dLon := transform(lon-105.0, lat-35.0)
	radLat := lat / 180.0 * math.Pi
	magic := math.Sin(radLat)
	magic = 1 - ee*magic*magic
	sMagic := math.Sqrt(magic)

	dLat = (dLat * 180.0) / ((a * (1 - ee)) / (magic * sMagic) * math.Pi)
	dLon = (dLon * 180.0) / (a / sMagic * math.Cos(radLat) * math.Pi)
	return lon + dLon, lat + dLat
}

// GCJ02toWGS84 火星坐标系->WGS84坐标系  精度为1m至2m
func GCJ02toWGS84(lon, lat float64) (float64, float64) {
	mgLon, mgLat := WGS84toGCJ02(lon, lat)
	return lon*2 - mgLon, lat*2 - mgLat
}

// GCJ02toWGS84Exact 火星坐标系->WGS84坐标系  精度小于0.5m 较GCJ02toWGS84慢15倍
func GCJ02toWGS84Exact(gcjLat, gcjLng float64) (wgsLat, wgsLng float64) {
	dLat, dLng := 0.01, 0.01
	mLat, mLng := gcjLat-dLat, gcjLng-dLng
	pLat, pLng := gcjLat+dLat, gcjLng+dLng
	for {
		wgsLat, wgsLng = (mLat+pLat)/2, (mLng+pLng)/2
		tmpLat, tmpLng := WGS84toGCJ02(wgsLat, wgsLng)
		dLat, dLng = tmpLat-gcjLat, tmpLng-gcjLng
		if math.Abs(dLat) < threshold && math.Abs(dLng) < threshold {
			return
		}
		if dLat > 0 {
			pLat = wgsLat
		} else {
			mLat = wgsLat
		}
		if dLng > 0 {
			pLng = wgsLng
		} else {
			mLng = wgsLng
		}
	}
}

// BD09toWGS84 百度坐标系->WGS84坐标系
func BD09toWGS84(lon, lat float64) (float64, float64) {
	return GCJ02toWGS84(BD09toGCJ02(lon, lat))
}

// WGS84toBD09 WGS84坐标系->百度坐标系
func WGS84toBD09(lon, lat float64) (float64, float64) {
	return GCJ02toBD09(WGS84toGCJ02(lon, lat))
}

// WebMCtoWGS84 球面墨卡托->WGS84坐标系
func WebMCtoWGS84(x, y float64) (float64, float64) {
	if !(x >= -mc && x <= mc) {
		return x, y
	}
	lng := x / mc * 180
	lat := y / mc * 180
	lat = 180 / math.Pi * (2*math.Atan(math.Exp(lat*math.Pi/180)) - math.Pi/2)
	return lng, lat
}

// WGS84toWebMC WGS84坐标系->球面墨卡托
func WGS84toWebMC(lon, lat float64) (float64, float64) {
	x := lon * mc / 180
	y := math.Log(math.Tan((90+lat)*math.Pi/360)) / (math.Pi / 180)
	y = y * mc / 180
	return x, y
}

func transform(x, y float64) (lat, lon float64) {
	absX := math.Sqrt(math.Abs(x))
	xPi, yPi := x*math.Pi, y*math.Pi
	d := 20.0*math.Sin(6.0*xPi) + 20.0*math.Sin(2.0*xPi)
	lat, lon = d, d
	lat += 20.0*math.Sin(yPi) + 40.0*math.Sin(yPi/3.0)
	lon += 20.0*math.Sin(xPi) + 40.0*math.Sin(xPi/3.0)
	lat += 160.0*math.Sin(yPi/12.0) + 320*math.Sin(yPi/30.0)
	lon += 150.0*math.Sin(xPi/12.0) + 300.0*math.Sin(xPi/30.0)
	lat *= 2.0 / 3.0
	lon *= 2.0 / 3.0
	lat += -100.0 + 2.0*x + 3.0*y + 0.2*y*y + 0.1*x*y + 0.2*absX
	lon += 300.0 + x + 2.0*y + 0.1*x*x + 0.1*x*y + 0.1*absX
	return
}

// 百度墨卡托
var mcBand = []float64{12890594.86, 8362377.87, 5591021, 3481989.83, 1678043.12, 0}
var mc2ll = [][]float64{
	[]float64{1.410526172116255e-8, 0.00000898305509648872, -1.9939833816331, 200.9824383106796, -187.2403703815547, 91.6087516669843, -23.38765649603339, 2.57121317296198, -0.03801003308653, 17337981.2},
	[]float64{-7.435856389565537e-9, 0.000008983055097726239, -0.78625201886289, 96.32687599759846, -1.85204757529826, -59.36935905485877, 47.40033549296737, -16.50741931063887, 2.28786674699375, 10260144.86},
	[]float64{-3.030883460898826e-8, 0.00000898305509983578, 0.30071316287616, 59.74293618442277, 7.357984074871, -25.38371002664745, 13.45380521110908, -3.29883767235584, 0.32710905363475, 6856817.37},
	[]float64{-1.981981304930552e-8, 0.000008983055099779535, 0.03278182852591, 40.31678527705744, 0.65659298677277, -4.44255534477492, 0.85341911805263, 0.12923347998204, -0.04625736007561, 4482777.06},
	[]float64{3.09191371068437e-9, 0.000008983055096812155, 0.00006995724062, 23.10934304144901, -0.00023663490511, -0.6321817810242, -0.00663494467273, 0.03430082397953, -0.00466043876332, 2555164.4},
	[]float64{2.890871144776878e-9, 0.000008983055095805407, -3.068298e-8, 7.47137025468032, -0.00000353937994, -0.02145144861037, -0.00001234426596, 0.00010322952773, -0.00000323890364, 826088.5},
}
var llBand = []float64{75, 60, 45, 30, 15, 0}
var ll2mc = [][]float64{
	[]float64{-0.0015702102444, 111320.7020616939, 1704480524535203, -10338987376042340, 26112667856603880, -35149669176653700, 26595700718403920, -10725012454188240, 1800819912950474, 82.5},
	[]float64{0.0008277824516172526, 111320.7020463578, 647795574.6671607, -4082003173.641316, 10774905663.51142, -15171875531.51559, 12053065338.62167, -5124939663.577472, 913311935.9512032, 67.5},
	[]float64{0.00337398766765, 111320.7020202162, 4481351.045890365, -23393751.19931662, 79682215.47186455, -115964993.2797253, 97236711.15602145, -43661946.33752821, 8477230.501135234, 52.5},
	[]float64{0.00220636496208, 111320.7020209128, 51751.86112841131, 3796837.749470245, 992013.7397791013, -1221952.21711287, 1340652.697009075, -620943.6990984312, 144416.9293806241, 37.5},
	[]float64{-0.0003441963504368392, 111320.7020576856, 278.2353980772752, 2485758.690035394, 6070.750963243378, 54821.18345352118, 9540.606633304236, -2710.55326746645, 1405.483844121726, 22.5},
	[]float64{-0.0003218135878613132, 111320.7020701615, 0.00369383431289, 823725.6402795718, 0.46104986909093, 2351.343141331292, 1.58060784298199, 8.77738589078284, 0.37238884252424, 7.45},
}

// BD09MCtoBD09 百度墨卡托->百度坐标系
func BD09MCtoBD09(x, y float64) (float64, float64) {
	x, y = math.Abs(x), math.Abs(y)
	var f []float64
	for k := range mcBand {
		if y >= mcBand[k] {
			f = mc2ll[k]
			break
		}
	}
	return convert(x, y, f)
}

// BD09MCtoWGS84 百度墨卡托->WGS84坐标系
func BD09MCtoWGS84(x, y float64) (float64, float64) {
	return BD09toWGS84(BD09MCtoBD09(x, y))
}

// BD09toBD09MC 百度坐标系->百度墨卡托
func BD09toBD09MC(lon, lat float64) (float64, float64) {
	lon = getLoop(lon, -180, 180)
	lat = getRange(lat, -74, 74)
	var f []float64
	for i := 0; i < len(llBand); i++ {
		if lat >= llBand[i] {
			f = ll2mc[i]
			break
		}
	}
	if len(f) > 0 {
		for i := len(llBand) - 1; i >= 0; i-- {
			if lat <= -llBand[i] {
				f = ll2mc[i]
				break
			}
		}
	}
	return convert(lon, lat, f)
}

func convert(x, y float64, f []float64) (float64, float64) {
	lon := f[0] + f[1]*math.Abs(x)
	cc := math.Abs(y) / f[9]

	var lat float64
	for i := 0; i <= 6; i++ {
		lat += f[i+2] * math.Pow(cc, float64(i))
	}

	if x < 0 {
		lon *= -1
	}
	if y < 0 {
		lat *= -1
	}
	return lon, lat
}

func getLoop(lng, min, max float64) float64 {
	for lng > max {
		lng -= max - min
	}
	for lng < min {
		lng += max - min
	}
	return lng
}

func getRange(lat, min, max float64) float64 {
	if min != 0 {
		lat = math.Max(lat, min)
	}
	if max != 0 {
		lat = math.Min(lat, max)
	}
	return lat
}
