package gobytes

import (
	"bytes"
	"fmt"
	"testing"
)

func TestSplitUsage(t *testing.T) {
	value := "michoi"
	fmt.Printf("%q\n", bytes.Split([]byte(value), []byte("")))
}

func TestSplitNUsage(t *testing.T) {
	config, err := LoadConfig("../config/app.conf")
	if err != nil {
		t.Fatal("err:", err.Error())
		return
	}
	fmt.Println("Config:", config)
}

func TestSplitNPrefix(t *testing.T) {
	value := "michoi"

	ns := []int{0, 1, 2, 3, 4, 100}
	for _, n := range ns {
		// if n is 1, result len is 1, r[0] = s
		val := bytes.SplitN([]byte(value), []byte("m"), n)
		if val == nil {
			fmt.Println("val is nil")
		} else {
			fmt.Printf("%d, %q\n", n, val)
		}
	}

	fmt.Printf("%q\n", bytes.SplitN([]byte("a,b,c"), []byte(","), 2))
	z := bytes.SplitN([]byte("a,b,c"), []byte(","), 0)
	fmt.Printf("%q (nil = %v)\n", z, z == nil)
}

func TestSplitNil(t *testing.T) {
	ns := []int{0, 1, 2, 3, 4, 100}
	for _, n := range ns {
		val := bytes.SplitN(nil, []byte(""), n)
		if val == nil {
			// n == 0
			fmt.Println("val is nil")
		} else {
			fmt.Printf("%q\n", val)
		}
	}
}

func TestSplitNSepIsNil(t *testing.T) {
	value := "michoi"

	ns := []int{0, 1, 2, 3, 4, 100}
	for _, n := range ns {
		val := bytes.SplitN([]byte(value), nil, n)
		if val == nil {
			fmt.Println("val is nil")
		} else {
			fmt.Printf("%q\n", val)
		}
	}
}

func TestSplitNNull(t *testing.T) {
	value := "michoi"

	ns := []int{-1, 0, 1, 2, 3, 4, 100}
	for _, n := range ns {
		val := bytes.SplitN([]byte(value), []byte(""), n)
		if val == nil {
			fmt.Println("val is nil")
		} else {
			fmt.Printf("%q\n", val)
		}
	}
}

func TestSplitNAllIsNil(t *testing.T) {
	ns := []int{0, 1, 2, 3, 4, 100}
	for _, n := range ns {
		val := bytes.SplitN(nil, nil, n)
		if val == nil {
			fmt.Println("val is nil")
		} else {
			fmt.Println(val)
		}
	}
}

func TestSplitAfterN(t *testing.T) {
	value := "michoi"

	ns := []int{0, 1, 2, 3, 4, 100}
	for n := range ns {
		val := bytes.SplitAfterN([]byte(value), []byte("i"), n)
		if val == nil {
			fmt.Println("val is nil")
		} else {
			fmt.Printf("%q\n", val)
		}
	}
}

func TestSplitAfter(t *testing.T) {
	fmt.Printf("%q\n", bytes.Split([]byte("a,b,c"), []byte(",")))
	fmt.Printf("%q\n", bytes.SplitAfter([]byte("a,b,c"), []byte(",")))
}
