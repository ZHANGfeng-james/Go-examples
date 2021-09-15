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



在初版的框架代码中，Engine 结构体定义中包含了**路由表**的实例，要知道路由表在整个 Web 框架中是很关键的一个实例，是否可以抽取出来形成独立的文件？这样也符号“单一职责原则”（在类型设计时）。接下来我们从 Engine 的源代码入手，抽取出 router.go 文件：

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

重构后，获得了一个单独的 router.go 文件，专门用于处理路由表相关的逻辑，符合“单一职责原则”。而且后续还可以在 router.go 中做更加重要的路由匹配策略（**动态路由**），让路由表的性能更加高效（**功能**、**性能**、**智能**）。

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

请求的构造是：请求行、请求头和请求体。对应响应也是这样的结构：响应行、响应头和响应体。构造一个完整的响应，需要考虑：StatusCode、Header 和 Body 部分，基本上每一次构造都需要考虑这些因素。如果不进行封装，那么框架的用户将需要写大量的冗余代码。

封装构造 http.ResponseWriter 的响应内容时，功能封装到哪里呢？是 Enginer 中？还是其他什么地方？此处，引入一个新的实体 Context（此 Context 和 context.Context 没有关系），将**每次请求**的 *http.Request 和 http.ResponseWriter 封装到 Context 类型实体中：

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

每一次 HTTP 的 Request 都会创建一个 Context 类型实例，而且符合 HTTP **和状态无关**的特征。因此，还需要重构 gee.go 和 router.go 文件：

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



