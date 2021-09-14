package gopath

import (
	"fmt"
	"path"
	"testing"
)

func TestBase(t *testing.T) {
	fmt.Println(path.Base("/a/b/"))
	fmt.Println(path.Base(""))
	fmt.Println(path.Base("//a////"))
}

func TestClean(t *testing.T) {
	paths := []string{
		"a/c/",
		"a//c/",
		"a/c/./",
		"a/c/b/..",
		"/../a/c",
		"/../a/b/../././/c",
		"a/b/../../xyz",
		"/",
		"",
	}

	for _, p := range paths {
		fmt.Printf("Clean(%q) = %q\n", p, path.Clean(p))
	}
}

func TestCleanOther(t *testing.T) {
	paths := []string{
		"a/b/../../../xyz",
		"a/b/../../xyz",
		"a/b/../xyz",
	}

	for _, p := range paths {
		fmt.Printf("Clean(%q) = %q\n", p, path.Clean(p))
	}
}

func TestDir(t *testing.T) {
	fmt.Println(path.Dir("/a/b/c"))
	fmt.Println(path.Dir("a/b/c"))
	fmt.Println(path.Dir("/a/"))
	fmt.Println(path.Dir("a/"))
	fmt.Println(path.Dir("/////////login/"))
	fmt.Println(path.Dir("/"))
	fmt.Println(path.Dir(""))
}

func TestJoin(t *testing.T) {
	fmt.Println(path.Join("a", "b", "c"))
	fmt.Println(path.Join("a", "b/c"))
	fmt.Println(path.Join("a/b", "c"))

	fmt.Println(path.Join("a/b", "../../../xyz"))

	fmt.Println(path.Join("", ""))
	fmt.Println(path.Join("a", ""))
	fmt.Println(path.Join("", "a"))
}

func TestPath(t *testing.T) {
	absolutePath := "/"
	relativePath := "/login"
	finalPath := path.Join(absolutePath, relativePath)

	fmt.Printf("finalPath:%s\n", finalPath)
}
