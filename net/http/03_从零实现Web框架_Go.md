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

接下来，把目光聚焦到 HTTP Request 上，让 Context 具备有**解析 URL 中参数**的能力：

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

net/http 标注库中的**路由匹配**核心逻辑是这样的：

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

下面把注意力放在：**如何实现前缀树**，考虑如下**疑惑**：

1. 前缀树是一个**树状结构**；
2. 随着注册路由的增多，前缀树的每一个节点下，还会**新增多个节点**；
3. 前缀树的节点可能会包含**一个 wild 字符**，此时这个节点称之为包含**通配符**的节点；
4. 作为一个前缀树，如果拿到了前缀树的 rootNode，就相当于可以遍历整个前缀树；
5. 需要为 HTTP 请求的每一种方法 GET、POST 等分别构建一棵前缀树；
6. 在路由匹配时，如何判断一个路由 path 已在前缀树中获得对应的匹配？

首先设计树节点上应该存储的信息量：

~~~go
// node constructor of router trie tree
type node struct {
	pattern  string  // 完整匹配路径
	part     string  // 当前节点的匹配内容
	children []*node // 每个节点下的子节点
	isWild   bool    // 是否包含通配符（* 和 :）
}
~~~

前缀树的构建和匹配，都是**一层一层**地经过**匹配**得到结果。path 和 node 匹配的逻辑：

~~~go
// matchChild matches children of node to find match one
func (n *node) matchChild(path string) *node {
	for _, ele := range n.children {
		if ele.part == path || ele.isWild {
			return ele
		}
	}
	return nil
}

// matchChildren matches all children, and return all nodes
func (n *node) matchChildren(path string) []*node {
	nodes := make([]*node, 0)
	for _, ele := range n.children {
		if ele.part == path || ele.isWild {
			nodes = append(nodes, ele)
		}
	}
	return nodes
}
~~~

比如前缀树中已注册了 `/p/:lang/doc` 的路由，在查找 `/p/go/doc` 时的过程是这样的：第一层精确匹配到了 `p`，第二层模糊匹配到了 `:lang`，对应设置参数 `[lang]=go`，再执行后续匹配。其中包含 `: / *` 通配符的 part 节点，其 isWild 值设置为 true。

作为一个前缀树——数据结构——其最关键的步骤就是**构建**和**查找**：

~~~go
// insert trie tree node with pattern
func (n *node) insert(pattern string, parts []string, height int) {
	//TEST CASE: /p/:name/join [p, :name, join] 0
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	// TDD
	// 0 --> p
	// 1 --> :name
	// 2 --> join

	//FIXME /p/:name/join /p/:time/sell
	//FIXME /p/:name /p/michoi
	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		// :name *filepath 存入 node.part
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

func (n *node) search(parts []string, height int) *node {
	// just for * only once
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			// middle path of route，并不是一个路由
			return nil
		}
		return n
	}

	part := parts[height]
	child := n.matchChildren(part)
	// children 为 []*node，是否存在多个匹配情况
	for _, item := range child {
		result := item.search(parts, height+1)
		if child != nil {
			return result
		}
	}
	return nil
}
~~~

在构建前缀树的过程中，递归查找每一层的节点，如果没有匹配到当前节点的 part，则新建一个节点。只有到 len(parts) 的最后才能为节点的 pattern 赋值当前的 pattern，中间所有创建的节点的 pattern 都应该设置为空字符串。因此，在路由匹配时，如果在最后的一次匹配中，发现节点的 pattern 为空字符，则说明路由前缀树中是没有注册该路由的。

路由的**基本数据结构**已经构建出来了，接下来需要将数据结构及其功能封装到 Router 中：

~~~go
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
~~~

作为一个 router，类型中包括了需要为每一种 HTTP method 构建的路由前缀树，roots 就是这些树的**数据结构**。比如 `roots["GET"]` 就是 GET 方法前缀树的根节点 `*node`。

向路由前缀树中**添加路由**：

~~~go
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
		router.roots[method] = &node{} // 每棵树的根节点都是空的
	}
	// insert(pattern string, parts []string, height int)
	router.roots[method].insert(pattern, parts, 0)

	key := method + "-" + pattern
	router.handlers[key] = handler
}
~~~

在路由前缀树中**查找** path 对应的 pattern：

~~~go
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
~~~

另外，在查询路由时，已经解析得到了对应的 URL 参数。在 Context 中增加了 Params 字段，用于保存参数：

~~~go
type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request

	Path   string
	Method string

	Params map[string]string

	StatusCode int
}

func (ctx *Context) Param(key string) (value string, ok bool) {
	value, ok = ctx.Params[key]
	return
}

func (router *router) handle(ctx *Context) {
	log.Printf("Receive: %4s - %s", ctx.Method, ctx.Path)
	if node, params := router.getRoute(ctx.Method, ctx.Path); node != nil {
		key := ctx.Method + "-" + node.pattern
		ctx.Params = params
		router.handlers[key](ctx)
	} else {
		ctx.String(http.StatusNotFound, "404 page not found, Path:%s", ctx.Path)
	}
}
~~~

对 router 这个模块功能做单元测试：

~~~go
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
~~~

单元测试中构造的路由前缀树：

![](./img/Snipaste_2021-09-17_16-37-46.png)

注册了 `/hello/:name` 和 `/hello/a/c`，下面来拿看看**路由匹配实例**：

1. HTTP GET `/hello/a`：对应的 pattern 是 `/hello/:name`
2. HTTP GET `/hello/a/c`：对应的 pattern 是 `/hello/a/c`，此时没有 name 参数

完成上面所有的内容后，下面看看用户如何使用：

~~~go
package main

import (
	"goweb/gee"
	"net/http"
)

func main() {
	engine := gee.New()

	...

	engine.GET("/hello/:name", func(ctx *gee.Context) {
		ctx.String(http.StatusOK, "hello %s, you're at %s\n", ctx.Param("name"), ctx.Path)
	})
	engine.GET("/hello", func(ctx *gee.Context) {
		ctx.String(http.StatusOK, "hello %s, you're at %s\n", ctx.Query("name"), ctx.Path)
	})
	engine.GET("/assets/*filepath", func(ctx *gee.Context) {
		ctx.JSON(http.StatusOK, gee.H{"filepath": ctx.Param("filepath")})
	})

	engine.Run(":9999")
}
~~~

### 路由分组

先来一波烧脑的疑惑：

1. 路由分组的概念是什么？
2. 为什么需要路由分组？
3. 路由分组和中间件有什么关系？
4. 如何实现路由分组？可以先想象一下路由分组应该是怎样实现的？如何让分组路由的使用尽可能简洁？



软件设计中的**概念**，比如某个模型等，很多都是来自现实的使用场景，也就是**需求**。

比如：路由分组，或分组路由。如果没有路由分组，我们就需要针对每一个路由分别进行控制。但是真实的业务场景种，往往**某一个组路由**需要**相似的处理**。例如：

* 以 `/post` 开头的路由**匿名**可访问。
* 以 `/admin` 开头的路由**需要鉴权**。
* 以 `/api` 开头的路由是 RESTful 接口，可以对接第三方平台，需要三方平台鉴权。

大部分情况下的路由分组，是以**相同的前缀**来区分的。因此，下面将要实现的分组控制也是以前缀来区分，并且支持分组的**嵌套**。另外还要考虑**中间件**在分组上的作用，比如 `/admin` 分组，可以应用鉴权中间件。

那考虑将路由分组做一个**模型**来实现，也就对应一个**类型**。下面就需要考虑这个类型需要**具备的属性**，另外 RouterGroup 对象还需要有访问 router 的能力，为了方便，可以在 RouterGroup 中，保存一个指针，指向 Engine：

~~~go
type (
	Engine struct {
		router       *router
		*RouterGroup // 内嵌*RouterGroup，*Engin类型具有*RouterGroup的所有方法
	}

	RouterGroup struct {
		engine     *Engine
		prefix     string
		parent     *RouterGroup // struct中不能定义相同类型的字段
		middleware []HandleFunc // middleware处理
	}
)
~~~

也就是将 Engine 作为了**最顶层的分组**，而且所有的 RouterGroup 都使用相同的 Engine 实例。

~~~go
func New() *Engine {
	engine := &Engine{router: newRouter()}
	// 为内嵌字段赋值的方法：engine.RouterGroup
	engine.RouterGroup = &RouterGroup{engine: engine}
	// engine.RouterGroup.prefix 为空
	// engine.RouterGroup.parent 为 nil
	// engine.RouterGroup.middleware 为 nil

	return engine
}

func (group *RouterGroup) NewRouterGroup(prefix string) *RouterGroup {
	// 所有的 RouterGroup 共享同一个*Engine
	engine := group.engine

	newGroup := &RouterGroup{}
	newGroup.engine = engine
	newGroup.prefix = group.prefix + prefix // 拼接 prefix
	newGroup.parent = group

	newGroup.middleware = nil

	return newGroup
}
~~~

紧接着将 Engine 的方法**重构**成 RouterGroup 的方法：

~~~go
func (group *RouterGroup) GET(path string, handler HandleFunc) {
	group.addRoute("GET", path, handler)
}

func (group *RouterGroup) POST(path string, handler HandleFunc) {
	group.addRoute("POST", path, handler)
}

func (group *RouterGroup) addRoute(method, component string, handler HandleFunc) {
	pattern := group.prefix + component // 拼接 group.prefixe 和 component
	log.Printf("component: %s, pattern: %s\n", component, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}
~~~

下面来看看用户的使用情况：

~~~go
package main

import (
	"goweb/gee"
	"net/http"
)

func main() {
	engine := gee.New()
	...
	helloGroup := engine.NewRouterGroup("/v1")
	{
		// 获得一个 RouterGroup 后，直接使用该类型注册 pattern
		helloGroup.GET("/:name", func(ctx *gee.Context) {
			ctx.String(http.StatusOK, "hello %s, you're at %s\n", ctx.Param("name"), ctx.Path)
		})
		helloGroup.GET("/geektutu/join", func(ctx *gee.Context) {
			ctx.String(http.StatusOK, "hello %s, you're at %s\n", ctx.Query("name"), ctx.Path)
		})
	}

	v2 := engine.NewRouterGroup("/v2")
	{
		v2.GET("/", func(ctx *gee.Context) {
			ctx.String(http.StatusOK, "you're at %s\n", ctx.Path)
		})
		v2.GET("/help", func(ctx *gee.Context) {
			ctx.String(http.StatusOK, "you're at %s\n", ctx.Path)
		})
	}

	engine.Run(":9999")
}
~~~

### 中间件

烧脑的疑惑：

1. 中间件是什么？为什么存在中间件这个概念？
2. 如何在 Web 框架中实现中间件？
3. 中间件如何触发执行？如何自定义中间件的执行顺序？
