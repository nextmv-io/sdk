package common_test

import (
	"testing"

	"github.com/nextmv-io/sdk/nextroute/common"
)

func BenchmarkDistance(b *testing.B) {
	d := common.NewDistance(100, common.Meters)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = d.Value(common.Meters)
	}
}
