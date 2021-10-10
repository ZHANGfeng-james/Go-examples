package gee

import (
	"fmt"
	"reflect"
	"testing"
)

func TestParsePattern(t *testing.T) {
	ok := reflect.DeepEqual(parsePattern("/p/:name"), []string{"p", ":name"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/*"), []string{"p", "*"})
	// parsePattern only for * just once
	ok = ok && reflect.DeepEqual(parsePattern("/p/*name/*"), []string{"p", "*name"})

	if !ok {
		t.Fatal("test parsePattern failed")
	}
}

func initTrieTree() *router {
	router := newRouter()

	// addRoute(method, pattern string, handler HandleFunc)
	router.addRoute("GET", "/", nil)
	router.addRoute("GET", "/hello/:name", nil)
	router.addRoute("GET", "/hello/b/c", nil)
	router.addRoute("GET", "/hi/:name", nil)
	router.addRoute("GET", "/assets/*filepath", nil)
	return router
}

func TestGetRoute(t *testing.T) {
	router := initTrieTree()

	path := "/hello/geektutu"
	// getRoute(method, path string) (*node, map[string]string)
	node, params := router.getRoute("GET", path)

	if node == nil {
		t.Fatal("there is a router for /hello/geektutu")
	}
	if node.pattern != "/hello/:name" {
		t.Fatal("pattern should be /hello/:name")
	}
	if params["name"] != "geektutu" {
		t.Fatal("param should be equal to 'geektutu'")
	}
	fmt.Printf("Path:%s, found: %s, params: %s\n", path, node.pattern, params["name"])
}

func TestGetRouteWithWildStar(t *testing.T) {
	router := initTrieTree()
	path := "/assets/file1.txt"
	node, params := router.getRoute("GET", path)
	ok := node.pattern == "/assets/*filepath"
	if !ok {
		t.Fatalf("Path: %s, pattern should be %s\n", path, "/assets/*filepath")
	}
	ok = params["filepath"] == "file1.txt"
	if !ok {
		t.Fatalf("Path:%s, params should be %s\n", path, "file1.txt")
	}
	fmt.Printf("Path:%s, found: %s, params: %s\n", path, node.pattern, params["filepath"])

	path = "/assets/dir/404.css"
	node, params = router.getRoute("GET", path)
	ok = node.pattern == "/assets/*filepath"
	if !ok {
		t.Fatalf("Path: %s, pattern should be %s\n", path, "/assets/*filepath")
	}
	ok = params["filepath"] == "dir/404.css"
	if !ok {
		t.Fatalf("Path:%s, params should be %s\n", path, "dir/404.css")
	}
	fmt.Printf("Path:%s, found: %s, params: %s\n", path, node.pattern, params["filepath"])
}
