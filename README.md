# transform
坐标系转换

- https://en.wikipedia.org/wiki/Restrictions_on_geographic_data_in_China
- https://github.com/sshuair/coord-convert
- https://github.com/googollee/eviltransform/blob/master/go/transform.go
- WGS84：即地球坐标系，国际上通用的坐标系。
- WebMC：即Web墨卡托，Web映射上的事实标准。
- ZXY：即瓦片坐标系，瓦片地图的切割方式。
- GCJ02：即火星坐标系，WGS84坐标系经加密后的坐标系。
- BD09：即百度坐标系，GCJ02坐标系经加密后的坐标系。
- BD09MC：即百度Web墨卡托，Web墨卡托经2次加密后的坐标系。

Installation
------

```bash
go get github.com/l-we/transform
```

Quick Start
------

```Go
package main

import (
	"fmt"

	"github.com/l-we/transform"
)

func main() {
	fmt.Println(transform.BD09toGCJ02(121.5272106, 31.1774276))
	fmt.Println(transform.GCJ02toBD09(121.5272106, 31.1774276))
	fmt.Println(transform.WGS84toGCJ02(121.5272106, 31.1774276))
	fmt.Println(transform.GCJ02toWGS84(121.5272106, 31.1774276))
	fmt.Println(transform.GCJ02toWGS84Exact(121.5272106, 31.1774276))
	fmt.Println(transform.BD09toWGS84(121.5272106, 31.1774276))
	fmt.Println(transform.WGS84toBD09(121.5272106, 31.1774276))

	fmt.Println(transform.WGS84toWebMC(121.5272106, 31.1774276))
	fmt.Println(transform.BD09toBD09MC(121.5272106, 31.1774276))

	fmt.Println(transform.WebMCtoWGS84(13529697, 3633994))
	fmt.Println(transform.BD09MCtoBD09(13529697, 3633994))
	fmt.Println(transform.BD09MCtoWGS84(13529697, 3633994))

	fmt.Println(transform.ZXYtoWebMCBound(12, 3429, 1673))
	fmt.Println(transform.ZXYtoWGS84Bound(12, 3429, 1673))
}
```