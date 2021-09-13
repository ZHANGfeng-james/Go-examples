package gostrings

import (
	"fmt"
	"strings"
	"testing"
)

func TestSplitFoo(t *testing.T) {
	values := "fof"
	fmt.Printf("%q\n", strings.Split(values, "f"))
}

func TestSplitNormal(t *testing.T) {
	values := "中-国-人"

	// 1 0
	fmt.Printf("%s\n", strings.Split(values, ""))
	// 1 1
	fmt.Printf("%s\n", strings.Split(values, "-"))
	fmt.Printf("%s\n", strings.Split("中国人", "-"))

	// 0 1
	fmt.Printf("%q\n", strings.Split("", "-"))

	// 0 0
	fmt.Printf("%q\n", strings.Split("", ""))

	result := strings.Split("", "-")
	fmt.Println(len(result))

	result = strings.Split("", "")
	fmt.Println(len(result))

	fmt.Printf("%q\n", strings.Split("a,b,c", ","))
	fmt.Printf("%q\n", strings.Split("a man a plan a canal panama", "a "))
	fmt.Printf("%q\n", strings.Split(" xyz ", ""))
	fmt.Printf("%q\n", strings.Split("", "Bernardo O'Higgins"))
}

func TestSplitAfter(t *testing.T) {
	values := "中-国-人"

	// 1 0
	fmt.Printf("%s\n", strings.SplitAfter(values, ""))
	// 1 1
	fmt.Printf("%s\n", strings.SplitAfter(values, "-"))
	fmt.Printf("%s\n", strings.SplitAfter("中国人", "-"))

	// 0 1
	fmt.Printf("%q\n", strings.SplitAfter("", "-"))

	// 0 0
	fmt.Printf("%q\n", strings.SplitAfter("", ""))

	result := strings.SplitAfter("", "-")
	fmt.Println(len(result))

	result = strings.SplitAfter("", "")
	fmt.Println(len(result))

	fmt.Printf("%q\n", strings.SplitAfter("a,b,c", ","))
	fmt.Printf("%q\n", strings.SplitAfter("a man a plan a canal panama", "a "))
	fmt.Printf("%q\n", strings.SplitAfter(" xyz ", ""))
	fmt.Printf("%q\n", strings.SplitAfter("", "Bernardo O'Higgins"))
}

func TestSplitN(t *testing.T) {
	values := "中-国-人"
	// 1 0
	fmt.Printf("%q\n", strings.SplitN(values, "", 10))
	// 1 1
	fmt.Printf("%q\n", strings.SplitN(values, "-", 10))
	fmt.Printf("%q\n", strings.SplitN("中国人", "-", 10))

	// 0 1
	fmt.Printf("%q\n", strings.SplitN("", "-", 10))

	// 0 0
	fmt.Printf("%q\n", strings.SplitN("", "", 10))

	fmt.Printf("%q\n", strings.SplitN("a,b,c", ",", 2))
	z := strings.SplitN("a,b,c", ",", 0)
	fmt.Printf("%q (nil = %v)\n", z, z == nil)
}

func TestSplitAfterN(t *testing.T) {
	values := "中-国-人"
	// 1 0
	fmt.Printf("%q\n", strings.SplitAfterN(values, "", 10))
	// 1 1
	fmt.Printf("%q\n", strings.SplitAfterN(values, "-", 10))
	fmt.Printf("%q\n", strings.SplitAfterN("中国人", "-", 10))

	// 0 1
	fmt.Printf("%q\n", strings.SplitAfterN("", "-", 10))

	// 0 0
	fmt.Printf("%q\n", strings.SplitAfterN("", "", 10))

	fmt.Printf("%q\n", strings.SplitAfterN("a,b,c", ",", 2))
	z := strings.SplitAfterN("a,b,c", ",", 0)
	fmt.Printf("%q (nil = %v)\n", z, z == nil)
}
