package common_test

import (
	"math/rand"
	"testing"

	"github.com/nextmv-io/sdk/nextroute/common"
)

func BenchmarkHaversine(b *testing.B) {
	r := rand.New(rand.NewSource(0))
	lon1, lat1 := r.Float64()*360-180, r.Float64()*180-90
	lon2, lat2 := r.Float64()*360-180, r.Float64()*180-90
	loc1, _ := common.NewLocation(lon1, lat1)
	loc2, _ := common.NewLocation(lon2, lat2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = common.Haversine(loc1, loc2)
	}
}
