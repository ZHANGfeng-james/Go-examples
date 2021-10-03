package singleflight

import "testing"

func TestSingleflight(t *testing.T) {
	group := &Group{}

	value, err := group.Do("key", func() (interface{}, error) {
		return "bar", nil
	})

	if err != nil || value != "bar" {
		t.Fatal("test singleflight failed")
	}
}
