package v1

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"testing"
)

type P struct {
	X, Y, Z int
	Name    string
}

type Q struct {
	X, Y *int32
	Name string
}

func TestEndingGob(t *testing.T) {
	var buf bytes.Buffer

	encoder := gob.NewEncoder(&buf)
	decoder := gob.NewDecoder(&buf)

	err := encoder.Encode(P{
		3, 4, 5,
		"Michoi",
	})
	if err != nil {
		log.Fatal("encode error:", err)
	}

	var instance Q
	err = decoder.Decode(&instance)
	if err != nil {
		log.Fatal("decode error:", err)
	}
	fmt.Printf("%q: {%d, %d}\n", instance.Name, *instance.X, *instance.Y)
}
