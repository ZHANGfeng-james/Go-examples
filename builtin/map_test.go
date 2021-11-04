package builtin

import (
	"math/rand"
	"testing"
	"time"
)

func TestMapCrate(t *testing.T) {
	createMap()
}

func TestMapCreateNil(t *testing.T) {
	createMapNil()
}

func TestMapSizeof(t *testing.T) {
	getMapSizeof()
}

func TestMapRetainValue(t *testing.T) {
	retainEle()
}

func TestMapCalledFunc(t *testing.T) {
	callFuncUseMap()
}

func TestMapForRange(t *testing.T) {
	forRangeMap()
}

func TestMapConcurrent(t *testing.T) {
	// mapConcurrent()
	mapGoAction()
}

func BenchmarkBuiltinMap(b *testing.B) {
	b.RunParallel(func(p *testing.PB) {
		r := rand.New(rand.NewSource(time.Now().Unix()))
		for p.Next() {
			k := r.Intn(100000000)
			storeUseBuiltinMap(k, k)
		}
	})
}

func TestMapKeyType(t *testing.T) {
	mapKey()
}

func TestMapForRangeAndDelete(t *testing.T) {
	forRangeAndDelete()
}

func TestMapForRangeAndAdd(t *testing.T) {
	forRangeAndAdd()
}
