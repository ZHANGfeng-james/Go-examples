package gostrings

import (
	"fmt"
	"strings"
	"testing"
)

func TestContains(t *testing.T) {
	// 1 1
	fmt.Println(strings.Contains("seafood", "foo"))
	fmt.Println(strings.Contains("seafood", "bar"))
	// 1 0
	fmt.Println(strings.Contains("seafood", ""))
	// 0 0
	fmt.Println(strings.Contains("", ""))
	// 0 1
	fmt.Println(strings.Contains("", "foo"))
}

func TestContainsAny(t *testing.T) {
	// 1 1
	fmt.Println(strings.ContainsAny("seafood", "foo"))
	fmt.Println(strings.ContainsAny("seafood", "bar"))
	// 1 0
	fmt.Println(strings.ContainsAny("seafood", ""))
	// 0 0
	fmt.Println(strings.ContainsAny("", ""))
	// 0 1
	fmt.Println(strings.ContainsAny("", "foo"))
}

func TestContainsRune(t *testing.T) {
	fmt.Println(strings.ContainsRune("众里寻他千百度，蓦然回首，那人却在灯火阑珊处！", '\u706b'))

	// Finds whether a string contains a particular Unicode code point.
	// The code point for the lowercase letter "a", for example, is 97.
	fmt.Println(strings.ContainsRune("aardvark", 97))
	fmt.Println(strings.ContainsRune("timeout", 97))
}
