package v4

import (
	"log"
	"strconv"
	"testing"
)

func TestHash(t *testing.T) {
	hash := New(3, func(data []byte) uint32 {
		// "06" --> 6; "16" --> 16
		i, _ := strconv.Atoi(string(data))
		return uint32(i)
	})

	hash.Add("6", "4", "2")

	testCases := map[string]string{
		"2":  "2",
		"11": "2",
		"23": "4",
		"27": "2",
	}

	for k, v := range testCases {
		got := hash.Get(k)
		log.Printf("ask for %s, got: %s", k, got)
		if got != v {
			t.Fatalf("Asking for %s, should have yielded %s", k, v)
		}
	}

	hash.Add("8")
	testCases["27"] = "8"
	testCases["29"] = "2"
	for k, v := range testCases {
		got := hash.Get(k)
		log.Printf("ask for %s, got: %s", k, got)
		if got != v {
			t.Fatalf("Asking for %s, should have yielded %s", k, v)
		}
	}
}
