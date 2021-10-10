package gee

import (
	"log"
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node      // roots key eg. roots["GET"] roots["POST"]
	handlers map[string]HandleFunc // handlers key eg. handlers["GET-/p/:name/join"] handlers["POST-/p/:name"]
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandleFunc),
	}
}

func parsePattern(pattern string) []string {
	// "/p/:name/join" --> [ ,p,:name,join] len() = 4
	parts := strings.Split(pattern, "/")
	vs := make([]string, 0)
	for _, item := range parts {
		if item != "" { // filte ""; for ""
			vs = append(vs, item)
			if item[0] == '*' { // only for * just once
				break // pattern:/p/*name/* --> [p, *name]
			}
		}
	}
	return vs
}

func (router *router) addRoute(method, pattern string, handler HandleFunc) {
	log.Printf("Route %4s - %s", method, pattern)
	if pattern == "" {
		panic("router pattern is empty path")
	}

	parts := parsePattern(pattern)

	if _, ok := router.roots[method]; !ok {
		router.roots[method] = &node{}
	}
	// insert(pattern string, parts []string, height int)
	router.roots[method].insert(pattern, parts, 0)

	key := method + "-" + pattern
	router.handlers[key] = handler
}

func (router *router) getRoute(method, path string) (*node, map[string]string) {
	root, ok := router.roots[method]
	if !ok {
		return nil, nil
	}

	// coding 锻炼写代码的逻辑，第一步做什么，第二步做什么...... Input/Output 分别是什么
	// read code 掌握代码背后的设计（思路和艺术），为什么这么设计，如果是我，我该如何设计
	parts := parsePattern(path)
	node := root.search(parts, 0)
	if node != nil {
		params := make(map[string]string)
		for index, item := range parsePattern(node.pattern) {
			if item[0] == ':' {
				params[item[1:]] = parts[index]
			}
			if item[0] == '*' && len(item) > 1 {
				params[item[1:]] = strings.Join(parts[index:], "/")
			}
		}
		return node, params
	}

	return nil, nil
}

func (router *router) handle(ctx *Context) {
	log.Printf("Receive: %4s - %s", ctx.Method, ctx.Path)
	if node, params := router.getRoute(ctx.Method, ctx.Path); node != nil {
		key := ctx.Method + "-" + node.pattern
		ctx.Params = params
		// 将注册的 HandleFunc 添加到 Context 中
		ctx.handlers = append(ctx.handlers, router.handlers[key])
	} else {
		// 添加异常处理的 HandleFunc
		ctx.handlers = append(ctx.handlers, func(ctx *Context) {
			ctx.String(http.StatusNotFound, "404 page not found, Path:%s", ctx.Path)
		})
	}
	ctx.Next()
}
