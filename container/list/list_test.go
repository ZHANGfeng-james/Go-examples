package list

import "testing"

func TestList(t *testing.T) {
	result := UseList()
	// []interface{}
	if value, ok := result.(string); ok {
		if value != "abcd" {
			t.Fatal("container/list run error")
		}
	}
}

func TestEmptyList(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Log("must use New to Create a List")
		}
	}()

	EmptyList()
}

func TestMarkIsNil(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Log("parameter in InsertXxx method should not be nil")
		}
	}()

	MarkIsNil()
}
