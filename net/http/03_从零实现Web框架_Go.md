> Web 框架，即 Web 服务端的脚手架，使用 net.http 包下的封装类实现，参考：gee-web，其最终都是以 Gin 为原型

Web 框架引发的**疑惑**：

1. 既然是是一个 Web 框架，那有哪几部分组成？
2. Gin、Beego 等 Web 框架的组成是怎样的？其基本功能有哪些
3. 一个 HTTP 请求经历了哪几个步骤？
4. 一个基本的后端服务程序的组织架构是怎样的？硬件层面包括哪些内容？涉及到哪些服务？这个部分可以和毛剑的极客时间课程联系起来。
5. HTTP 请求的方式、方法有哪些？



### Web 框架雏形

先来一波**烧脑的思考题**：

1. 如何创建一个 Server，并接收 Client 的 HTTP request？

   解答：参考 [02_使用3种思路构建Web服务器](./02_使用3种思路构建Web服务器.md)
   
2. Server 接收到 HTTP 请求后，**经历了哪些步骤**后 Request 得到了处理？比如如何**匹配** URL 的？

3. Server 是如何**并发**接收大量的 Client 端 HTTP Request 的？如何测量 Server 的 QPS？



Go 语言标准库 net/http 封装了 **HTTP 网络编程**的基础接口，该包可用于封装一个 Web 框架（**Gin 使用的就是该标准库**）。在不知道 Gin 怎样实现 Web 框架 时，认为 Gin 肯定做了大量的工作，但实际上仍然是在标准库的基础上实现的。

**如何创建一个 Server，并接收 Client 的 HTTP request**？为了探求这个问题的答案，参考标准库中的 net/http 包的使用：

~~~go
package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", indexHandleFunc)
	http.HandleFunc("/hello", helloHandlFunc)

	log.Fatal(http.ListenAndServe(":9999", nil))
}

func indexHandleFunc(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, req.RequestURI)
}

func helloHandlFunc(w http.ResponseWriter, req *http.Request) {
	for key, value := range req.Header { // req.Header map[string][]string
		fmt.Fprintf(w, "key:%s, value:%s\n", key, value)
	}
}
~~~

`http.ListenAndServe` 实现 HTTP 服务器的启动，同时持续接收 Client 的 HTTP 请求。

Server 端接收到 HTTP 请求后（Accept），创建 goroutine 并对这个 Request 做具体的处理。此处“具体的处理”就包括依据 HTTP Request 的 URI 在 DefaultServeMux 中查找已注册的 path（**路由**） 和 handler，若找到，则执行对应的 handler，否则直接返回：`404 page not found`

~~~go
// HandleFunc registers the handler function for the given pattern
// in the DefaultServeMux.
// The documentation for ServeMux explains how patterns are matched.
func HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
	DefaultServeMux.HandleFunc(pattern, handler)
}
~~~

对于 `/` 这个 path，可以匹配很多个 URI，比如：`/hello`（如果没有注册对应的 `/hello` 的 HandleFunc 时） 等。

~~~shell
C:\Users\Administrator>curl http://localhost:9999/hello
User-AgentAccept
C:\Users\Administrator>curl http://localhost:9999/hello
key:User-Agent, value:[curl/7.55.1]
key:Accept, value:[*/*]

C:\Users\Administrator>curl http://localhost:9999/helloo
/helloo
~~~

main 函数的最后 `http.ListenAndServer` 执行时，启动了 Web 服务，监听的端口是 `:9999`，第二个参数表示处理所有 HTTP 请求的 Handler 实例。

~~~go
// ListenAndServe listens on the TCP network address addr and then calls
// Serve with handler to handle requests on incoming connections.
// Accepted connections are configured to enable TCP keep-alives.
//
// The handler is typically nil, in which case the DefaultServeMux is used.
//
// ListenAndServe always returns a non-nil error.
func ListenAndServe(addr string, handler Handler) error {
	server := &Server{Addr: addr, Handler: handler}
	return server.ListenAndServe()
}

type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}
~~~

第二个参数一般是 nil，默认使用的是 DefaultServeMux 实例。

这第二个参数是**基于 net/http 标准库实现 Web 框架的入口**！只要传入实现了 Handler 这个接口的实例，所有的 HTTP Request 都会被该实例处理。

那接下来实现这个接口定制 HTTP Request 的处理：

~~~go
package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	engine := &Engine{}
	log.Fatal(http.ListenAndServe(":9999", engine))
}

type Engine struct{}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path // 获取到对应的URL，并依据URL匹配到HandleFunc
	switch path {
	case "/": // 内容就是HandleFunc
		fmt.Fprintf(w, req.RequestURI)
	case "/hello":
		for key, value := range req.Header { // req.Header map[string][]string
			fmt.Fprintf(w, "key:%s, value:%s\n", key, value)
		}
	default:
		fmt.Fprintf(w, "404 page not found")
	}
}
~~~

在 `ServeHTTP(w http.ResponseWriter, req *http.Request)` 这个方法声明中，`http.Request` 类型包含了本次 HTTP Request 的所有信息，包括请求地址、Header 和 Body 等。而 `http.ResponseWriter` 用于构造请求的响应。

在传入 engine 后，我们实现了将所有 HTTP 请求转向了我们自己的处理逻辑：

~~~go
func (sh serverHandler) ServeHTTP(rw ResponseWriter, req *Request) {
    // 此处可以自定义 handler，如果为空，则使用 DefaultServeMux
	handler := sh.srv.Handler
	if handler == nil {
		handler = DefaultServeMux
	}
	if req.RequestURI == "*" && req.Method == "OPTIONS" {
		handler = globalOptionsHandler{}
	}
    // 调用逻辑❼
    // 获取 ServeMux，并处理 Request
	handler.ServeHTTP(rw, req)
}
~~~

这个步骤的目的在于：拦截了所有的 HTTP Request，让其拥有的**统一的入口**。也就是说，HTTP Request 最先会到达这个方法，我可以在此处自定义路由映射规则，添加一些日志、校验、拦截等逻辑。

那接下来构造这个框架，这个框架以 Engine 为主功能类型，需要实现的功能有：

* 创建 Engine 实例
* 统一 HTTP Request 的入口
* 可添加 path - HandleRequest 处理，并实现根据 URL 匹配到对应的 HandleFunc

实现的框架如下：

~~~go
package gee

import (
	"fmt"
	"log"
	"net/http"
)

type Engine struct {
	router map[string]HandleFunc
}

func New() *Engine {
	return &Engine{router: make(map[string]HandleFunc)}
}

func (engine *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) GET(path string, handler HandleFunc) {
	engine.addRoute("GET", path, handler)
}

func (engine *Engine) POST(path string, handler HandleFunc) {
	engine.addRoute("POST", path, handler)
}

func (engine *Engine) addRoute(method, pattern string, handler HandleFunc) {
	key := method + "-" + pattern
	log.Printf("Route %4s - %s", method, pattern)
	engine.router[key] = handler
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Printf("Receive: %4s - %s", req.Method, req.URL.Path) // 不需要在 format 的结尾增加 newline

	key := req.Method + "-" + req.URL.Path
	if handler, ok := engine.router[key]; ok {
		handler(w, req)
	} else {
		fmt.Fprintf(w, "404 page not found")
	}
}

type HandleFunc func(w http.ResponseWriter, req *http.Request)

~~~

使用 goweb 框架后，实现 Server 代码变得精简了：

~~~go
package main

import (
	"fmt"
	"goweb/gee"
	"net/http"
)

func main() {
	engine := gee.New()

	engine.GET("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, req.RequestURI)
	})

	engine.GET("/hello", func(w http.ResponseWriter, req *http.Request) {
		for key, value := range req.Header { // req.Header map[string][]string
			fmt.Fprintf(w, "key:%s, value:%s\n", key, value)
		}
	})

	engine.Run(":9999")
}
~~~

重点来分析**首个框架版本**的实现：

1. HandleFunc 类型来自 net/http 标准库，让用户定义**路由**映射的**处理方法**（HandleFunc 封装）。
2. Engine 中封装了一个**路由表**，其类型是 `map[string]HandlFunc` 类型的值，其中 string 的构成是：`Method-URL`，比如 `GET-/`、`GET-/hello` 等。针对相同的路由（URL），如果请求方法不同，可以有不同的处理方法。
3. Engine 中实现的 `ServeHTTP(w http.ResponseWriter, req *http.Request)` 的作用就是解析出请求的 URL，并在路由表中查找对应的 HandleFunc，若找到则处理，反之则反馈 `404 page not found`。

测试用例：

~~~shell
C:\Users\Administrator>curl http://localhost:9999/
/
C:\Users\Administrator>curl http://localhost:9999/hello
key:User-Agent, value:[curl/7.55.1]
key:Accept, value:[*/*]

C:\Users\Administrator>curl http://localhost:9999/world
404 page not found
~~~

到目前为止， Go Web 框架的原型已经出来的，但是上述框架代码**并没有实现比标准库更强大的能力**（但仍实现了**路由注册表**，提供了用户**注册静态路由**的方法，包装了启动服务的函数），这就是我们后续要做的事情。

### 设计 Context 上下文

先来一波**烧脑的思考题**：

1. HTTP 的 Response 如何构建，是如何反馈给 Client 端的？
2. Response 的类型有哪些？JSON、HTML 等。作为一个 Web 框架，是否可以使用提取公共的代码组成方法，方便用户使用？
3. HTTP Request 中的参数是如何被解析的？



在初版的框架代码中，Engine 结构体定义中包含了**路由表**的实例，要知道路由表在整个 Web 框架中是很关键的一个实例，是否可以抽取出来形成独立的文件？这样也符合**“单一职责原则”（在类型设计时）**。接下来我们从 Engine 的源代码入手，抽取出 router.go 文件：

~~~go
package gee

import (
	"fmt"
	"log"
	"net/http"
)

type router struct {
	handlers map[string]HandleFunc
}

func newRouter() *router {
	return &router{handlers: make(map[string]HandleFunc)}
}

func (router *router) addRoute(method, pattern string, handler HandleFunc) {
	log.Printf("Route %4s - %s", method, pattern)
	key := method + "-" + pattern
	router.handlers[key] = handler
}

func (router *router) handle(w http.ResponseWriter, req *http.Request) {
	log.Printf("Receive: %4s - %s", req.Method, req.URL.Path)
	key := req.Method + "-" + req.URL.Path
	if handler, ok := router.handlers[key]; ok {
		handler(w, req)
	} else {
		fmt.Fprintf(w, "404 page not found")
	}
}
~~~

抽取出了 router.go 后，还需要对使用 router 的地方做**重构**：

~~~go
package gee

import (
	"net/http"
)

type Engine struct {
	router *router // 此处并没有内嵌 router 结构体，而使用了 *router 字段承接一个实例
}

func New() *Engine {
	return &Engine{router: newRouter()} // 重构
}

func (engine *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) GET(path string, handler HandleFunc) {
	engine.addRoute("GET", path, handler)
}

func (engine *Engine) POST(path string, handler HandleFunc) {
	engine.addRoute("POST", path, handler)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	engine.router.handle(w, req) // 重构
}

func (engine *Engine) addRoute(method, pattern string, handler HandleFunc) {
	engine.router.addRoute(method, pattern, handler) // 重构
}

type HandleFunc func(w http.ResponseWriter, req *http.Request)
~~~

重构后，获得了一个单独的 router.go 文件，专门用于处理路由表相关的逻辑，符合**“单一职责原则”**。而且后续还可以在 router.go 中做更加重要的路由匹配策略（**动态路由**），让路由表的性能更加高效（**功能**、**性能**、**智能**）。

接下来，将目光聚焦到 Response 上：Server 接收到 *http.Request，经过一系列的处理，最终总是需要反馈给 Client 消息的，也就是构造响应 http.ResponseWriter。而上述这两个标准库封装的功能粒度太细，用户在使用时难免会感受到繁琐，比如：

~~~go
package main

import (
	"encoding/json"
	"fmt"
	"goweb/gee"
	"net/http"
)

func main() {
	engine := gee.New()

	engine.GET("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, req.RequestURI)
	})

	engine.GET("/hello", func(w http.ResponseWriter, req *http.Request) {
		obj := map[string]interface{}{
			"name":     "geektutu",
			"password": 1234,
		}
        // 输出到 http.ResponseWriter 的流程和 Error 中的类似
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(obj); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	engine.Run(":9999")
}

// Error replies to the request with the specified error message and HTTP code.
// It does not otherwise end the request; the caller should ensure no further
// writes are done to w.
// The error message should be plain text.
func Error(w ResponseWriter, error string, code int) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	fmt.Fprintln(w, error)
}
~~~

请求的构造是：**请求行、请求头和请求体**。对应响应也是这样的结构：响应行、响应头和响应体。构造一个完整的响应，需要考虑：StatusCode、Header 和 Body 部分，基本上每一次构造都需要考虑这些因素。如果不进行封装，那么框架的用户将需要写大量的冗余代码。

封装构造 http.ResponseWriter 的响应内容时，**功能封装到哪里呢？**是 Enginer 中？还是其他什么地方？此处，引入一个新的实体 Context（此 Context 和 context.Context 没有关系），将**每次请求**的 *http.Request 和 http.ResponseWriter 封装到 Context 类型实体中：

~~~go
package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request

	Path   string
	Method string

	StatusCode int
}

func (ctx *Context) String(statusCode int, format string, values ...interface{}) {
	ctx.SetHeader("Content-Type", "text/plain")
	ctx.SetStatus(statusCode)
	ctx.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (ctx *Context) JSON(statusCode int, obj interface{}) {
	ctx.SetHeader("Content-Type", "application/json")
	ctx.SetStatus(statusCode)
	encoder := json.NewEncoder(ctx.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(ctx.Writer, err.Error(), http.StatusInternalServerError)
	}
}

func (ctx *Context) HTML(statusCode int, html string) {
	ctx.SetHeader("Content-Type", "text/html")
	ctx.SetStatus(statusCode)
	ctx.Writer.Write([]byte(html))
}

func (ctx *Context) Data(statusCode int, data []byte) {
	ctx.SetStatus(statusCode)
	ctx.Writer.Write(data)
}

func (ctx *Context) SetHeader(key, value string) {
	ctx.Writer.Header().Set(key, value)
}

func (ctx *Context) SetStatus(statusCode int) {
    ctx.StatusCode = statusCode
	ctx.Writer.WriteHeader(statusCode)
}
~~~

我们构建起来了 Context 结构体类型，并在其上创建了对应的方法：封装了以 JSON、XML、String、Data 格式输出的 http.ResponseWriter 方法，方便用户直接调用。

**每一次** HTTP 的 Request 都会**创建一个 Context 类型实例**，而且符合 HTTP **和状态无关**的特征。因此，还需要重构 gee.go 和 router.go 文件：

~~~go
package gee

import (
	"net/http"
)

type HandleFunc func(ctx *Context) // 重构

type Engine struct {
	router *router
}

func New() *Engine {
	return &Engine{router: newRouter()}
}

func (engine *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) GET(path string, handler HandleFunc) {
	engine.addRoute("GET", path, handler)
}

func (engine *Engine) POST(path string, handler HandleFunc) {
	engine.addRoute("POST", path, handler)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := newContext(w, req) // 重构，在此处创建 Context 类型实体，每一次都是一个和当次 Request 直接相关的 Context 类型实例
	engine.router.handle(ctx) // 重构
}

func (engine *Engine) addRoute(method, pattern string, handler HandleFunc) {
	engine.router.addRoute(method, pattern, handler)
}
~~~

以及更重要的：

~~~go
package gee

import (
	"fmt"
	"log"
)

type router struct {
	handlers map[string]HandleFunc
}

func newRouter() *router {
	return &router{handlers: make(map[string]HandleFunc)}
}

func (router *router) addRoute(method, pattern string, handler HandleFunc) {
	log.Printf("Route %4s - %s", method, pattern)
	key := method + "-" + pattern
	router.handlers[key] = handler
}

func (router *router) handle(ctx *Context) { // 重构
	log.Printf("Receive: %4s - %s", ctx.Method, ctx.Path)
	key := ctx.Method + "-" + ctx.Path
	if handler, ok := router.handlers[key]; ok {
		handler(ctx) // 重构
	} else {
		fmt.Fprintf(ctx.Writer, "404 page not found")
	}
}
~~~

接下来，把目光聚焦到 HTTP Request 上，让 Context 具备有解析 URL 中参数的能力：

~~~go
func (ctx *Context) postForm(key string) string {
	return ctx.Request.FormValue(key)
}

func (ctx *Context) Query(key string) string {
	return ctx.Request.URL.Query().Get(key) // Query是从URL中查询
}
~~~

另外还在 context.go 中新增一个**自定义类型**：

~~~go
type H map[string]interface{}
~~~

现在来看，整个代码已经很**清晰**了，模块的组成部分**各司其职**。

下面看看封装后，**框架的应用**情况：

~~~go
package main

import (
	"goweb/gee"
	"net/http"
)

func main() {
	engine := gee.New()

	engine.GET("/", func(ctx *gee.Context) {
		ctx.HTML(http.StatusOK, "<h1>Hello Gee<h1>")
	})

	engine.GET("/json", func(ctx *gee.Context) {
		obj := gee.H{
			"name":     "geektutu",
			"password": 1234,
		}
		ctx.JSON(http.StatusOK, obj)
	})

	engine.POST("/postform", func(ctx *gee.Context) { // 必须是 POST 请求，才能解析出 PostForm 内容
		ctx.JSON(http.StatusOK, gee.H{
			"name":     ctx.PostForm("name"),
			"password": ctx.PostForm("password"),
		})
		// example: curl "http://localhost:9999/postform" -X POST -d 'password=1&name=1'
	})

	engine.GET("/query", func(ctx *gee.Context) {
		username := ctx.Query("username")
		ctx.String(http.StatusOK, "Hello, %s!", username)
		// example: curl "http://localhost:9999/query?username=Michoi"
	})

	engine.Run(":9999")
}
~~~

特别注意，Windows 平台上使用 cmd 做 curl 网络请求：

~~~shell
curl "http://localhost:9999/postform" -X POST -d 'password=1&name=1'
~~~

**执行异常**，无法得到正确的请求结果！但是在 git 终端却**工作正常**。

### 路由表 Router

下来一波烧脑的疑惑：

1. 标准库 net/http 中路由表是如何创建的，如何匹配路由获得对应的 HandlerFunc？
2. 路由表是否可自定义，以此获得**更高的路由查找效率**？
3. 如何去实现**动态路由**？比如去实现既能匹配 `/hello/a` 也能匹配 `/hello/b` 的路由。



先来解答一个疑惑：为什么注册了 `/` 的 GET HandlerFunc，但是如果请求的是 `/anything` 时，对应执行了 `/` 的 HandleFunc？

这个疑惑涉及到 net/http 中**路由表的创建**，以及对应**路由匹配的逻辑**，也就是分为上面 2 个部分。解答如下：

~~~go
type ServeMux struct {
	mu    sync.RWMutex
	m     map[string]muxEntry
	es    []muxEntry // slice of entries sorted from longest to shortest.
	hosts bool       // whether any patterns contain hostnames
}

type muxEntry struct {
	h       Handler
	pattern string
}

// Handle registers the handler for the given pattern.
// If a handler already exists for pattern, Handle panics.
func (mux *ServeMux) Handle(pattern string, handler Handler) { // 路由表的创建
	mux.mu.Lock()
	defer mux.mu.Unlock()

	if pattern == "" {
		panic("http: invalid pattern")
	}
	if handler == nil {
		panic("http: nil handler")
	}
	if _, exist := mux.m[pattern]; exist { // 进入到 ServeMux 的 path，对应就是pattern
		panic("http: multiple registrations for " + pattern)
	}

	if mux.m == nil {
		mux.m = make(map[string]muxEntry)
	}
    e := muxEntry{h: handler, pattern: pattern} 
    mux.m[pattern] = e // mux.m: pattern - (pattern, handler) 的map结构
	if pattern[len(pattern)-1] == '/' { // 特殊的，以 '/' 结尾的pattern，添加到 mux.es 中
		mux.es = appendSorted(mux.es, e)
	}

	if pattern[0] != '/' {
		mux.hosts = true
	}
}

func appendSorted(es []muxEntry, e muxEntry) []muxEntry {
	n := len(es)
	i := sort.Search(n, func(i int) bool {
        // slice of entries sorted from longest to shortest.
		return len(es[i].pattern) < len(e.pattern)
	})
	if i == n {
		return append(es, e)
	}
    
	// we now know that i points at where we want to insert
	es = append(es, muxEntry{}) // try to grow the slice in place, any entry works.
	copy(es[i+1:], es[i:])      // Move shorter entries down
	es[i] = e
	return es
}
~~~

匹配路由的过程是这样的：

~~~go
// ServeHTTP dispatches the request to the handler whose
// pattern most closely matches the request URL.
func (mux *ServeMux) ServeHTTP(w ResponseWriter, r *Request) {
	if r.RequestURI == "*" {
		if r.ProtoAtLeast(1, 1) {
			w.Header().Set("Connection", "close")
		}
		w.WriteHeader(StatusBadRequest)
		return
	}
	h, _ := mux.Handler(r) // 寻找 HandlerFunc，相当于是入口
	h.ServeHTTP(w, r)
}

func (mux *ServeMux) Handler(r *Request) (h Handler, pattern string) {
    ...
	return mux.handler(host, r.URL.Path)
}

func (mux *ServeMux) handler(host, path string) (h Handler, pattern string) {
	mux.mu.RLock()
	defer mux.mu.RUnlock()

	// Host-specific pattern takes precedence over generic ones
	if mux.hosts {
		h, pattern = mux.match(host + path)
	}
	if h == nil {
		h, pattern = mux.match(path)
	}
	if h == nil {
		h, pattern = NotFoundHandler(), ""
	}
	return
}

// Find a handler on a handler map given a path string.
// Most-specific (longest) pattern wins.
func (mux *ServeMux) match(path string) (h Handler, pattern string) {
	// Check for exact match first. 在 map[string]muxEntry 中作精确查找
	v, ok := mux.m[path]
	if ok {
		return v.h, v.pattern
	}

	// Check for longest valid match.  mux.es contains all patterns
	// that end in / sorted from longest to shortest.
	for _, e := range mux.es {
		if strings.HasPrefix(path, e.pattern) { // 判断 path 是否具有 e.pattern 的前缀
			return e.h, e.pattern
		}
	}
	return nil, ""
}
~~~

路由匹配的核心逻辑是这样的：

1. 在 `mux.m[path]` 中作 path 的精确匹配，若查找到则返回；
2. 紧接着，在 `mux.es` 中查找，其排列顺序是 `mux.es[index].pattern` 从长到短依次排列的，且 pattern 的结尾是 `/` 字符。如果 path 具有某个 pattern 前缀，则表示匹配上了，并返回 handler；
3. 否则，返回 `404 page not found`

> 在注册阶段，path 进到进到标准库中是 pattern；在路由匹配阶段，path 就是路由——`r.URL.Path`。

举个例子：

~~~go
func main() {
	http.HandleFunc("/", indexHandleFunc)

	http.HandleFunc("/hel/", helloHandleFunc)

	log.Fatal(http.ListenAndServe(":9999", nil))
}
~~~

当请求：`curl http://localhost:9999/hel/hello` 会执行 helloHandleFunc，请求 `curl http://localhost:9999/hello` 会访问 indexHandleFunc。

另外，在上面标准库 net/http 的路由匹配时，是一种**静态路由匹配**，也就是说必须已注册的 pattern-handler。那如何实现**更灵活**的**动态路由**呢？

> 所谓**动态路由**，即一条路由规则可以**匹配某一类型**而**非某一条固定的路由**。例如`/hello/:name`，可以匹配`/hello/geektutu`、`hello/jack`等。

那接下来，要去**创建更高效率的路由匹配策略**，其中就包括寻找**适合当前问题场景**的**数据结构**：

动态路由有很多种实现方式，支持的规则、性能等有很大的差异。实现动态路由最常用的数据结构，被称为**前缀树**(Trie树)。看到名字你大概也能知道前缀树长啥样了：**每一个节点的所有的子节点都拥有相同的前缀**。这种结构非常适用于路由匹配：

- /:lang/doc
- /:lang/tutorial
- /:lang/intro
- /about
- /p/blog
- /p/related

HTTP请求的路径恰好是**由`/`分隔的多段**构成的，因此，**每一段**可以作为**前缀树的一个节点**。我们通过树结构查询，如果**中间某一层的节点**都**不满足条件**，那么就说明**没有匹配到的路由**，查询结束。

![](./img/trie_router.jpg)

创建的动态路由需具备如下功能：

* **参数匹配`:`**。例如 `/p/:lang/doc`，可以匹配 `/p/c/doc` 和 `/p/go/doc`。
* **通配`*`**。例如 `/static/*filepath`，可以匹配`/static/fav.ico`，也可以匹配`/static/js/jQuery.js`，这种模式常用于**静态服务器**，能够**递归**地匹配子路径。

