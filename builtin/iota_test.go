package builtin

import (
	"log"
	"testing"
)

func TestIotaCreate(t *testing.T) {
	var size Size = extraLarge
	log.Println(size)
}

func TestIotaValue(t *testing.T) {
	if c != 102 {
		t.Fail()
	}
}
