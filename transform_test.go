package transform

import (
	"testing"
)

func Test_transform(t *testing.T) {
	t.Log(BD09toBD09MC(121.431863, 31.027647))
}

func Benchmark_transform(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GCJ02toWGS84Exact(121.431863, 31.027647)
	}
}
