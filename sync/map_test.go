package sync

import (
	"math/rand"
	"testing"
	"time"
)

func TestSyncMap(t *testing.T) {
	syncMap()
}

func BenchmarkSyncMap(b *testing.B) {
	b.RunParallel(func(p *testing.PB) {
		r := rand.New(rand.NewSource(time.Now().Unix()))
		for p.Next() {
			k := r.Intn(100000000)
			storeUseSyncMap(k, k)
		}
	})
}
