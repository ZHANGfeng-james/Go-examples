package v1

import "testing"

func TestChangeValue(t *testing.T) {
	value := 10
	changeValue(&value)
	if value != 20 {
		t.Fatal("value is not changed!")
	}
}

func TestCreateNewIntValue(t *testing.T) {
	createNewIntValue()
}
