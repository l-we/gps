package transform

import "testing"

func Test_Num2WebMCBound(t *testing.T) {
	t.Log(ZXYtoWebMCBound(12, 3429, 1673))
	t.Log(ZXYtoWGS84Bound(12, 3429, 1673))
}
