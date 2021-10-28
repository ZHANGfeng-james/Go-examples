package builtin

import "testing"

func TestMapCrate(t *testing.T) {
	createMap()
}

func TestMapCreateNil(t *testing.T) {
	createMapNil()
}

func TestMapSizeof(t *testing.T) {
	getMapSizeof()
}
