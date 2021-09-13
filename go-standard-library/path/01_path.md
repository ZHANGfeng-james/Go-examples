Go 中的 path 标准库，是作为处理**斜线分隔的路径**的工具类。

path 标准库仅能用于处理**斜线分隔**的路径，比如 **URL 路径**。相应的，path 标准库不能用于处理 Windows 路径，比如：带有磁盘驱动符号以及反斜线。如果是要处理类似操作系统路径的问题，可以使用 path/filepath 标准库。

在整个 path 标准库中，有如下的疑惑：

1. path 到底表示什么？为什么是要应用在 URLs 场景种？
2. `/` 为什么表示根路径？



# 1 Base

`func Base(path string) string`：返回 path 的**最后一个元素**。在截取最后一个元素前，会删除最后的斜线（trailing slash）。如果入参路径为空，则返回 `.`；如果路径全部都是斜线，则返回 `/`。

~~~go
func TestBase(t *testing.T) {
	fmt.Println(path.Base("/a/b"))
	fmt.Println(path.Base(""))
	fmt.Println(path.Base("///"))
}
b
.
/
~~~

# 2 Clean

`func Clean(path string) string`：正如 Clean 的含义，函数返回**和入参 path 语义相同的最短路径名**。其处理依次遵循如下规则：

1. 用单个斜线替换多个连续的斜线；
2. 删除每一个 `.` 路径名（表示当前目录）；
3. 删除所有内部的 `..` 路径名，以及在这个 `..` 路径名之前的非 `..` 元素；
4. 使用 `/` 替换 `/..`，注意此种替换必须是在 path 的开头部分。

如果入参路径是 `/`，那返回值就是 `/`；如果经过 Clean 的处理，其返回值是空字符串，那返回值会是 `.`。

比如：

~~~go
func TestClean(t *testing.T) {
	paths := []string{
		"a/c",
		"a//c",
		"a/c/.",
		"a/c/b/..",
		"/../a/c",
		"/../a/b/../././/c",
        "a/b/../../../xyz",
		"/",
		"",
	}

	for _, p := range paths {
		fmt.Printf("Clean(%q) = %q\n", p, path.Clean(p))
	}
}

Clean("a/c") = "a/c"
Clean("a//c") = "a/c"
Clean("a/c/.") = "a/c"
Clean("a/c/b/..") = "a/c"
Clean("/../a/c") = "/a/c"
Clean("/../a/b/../././/c") = "/a/c"
Clean("a/b/../../../xyz") = "../xyz"
Clean("/") = "/"
Clean("") = "."
~~~

这个函数，会**删除 path 末尾的斜线**。

特别有意思的处理是：

~~~go
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

Clean("a/b/../../../xyz") = "../xyz"
Clean("a/b/../../xyz") = "xyz"
Clean("a/b/../xyz") = "a/xyz"
~~~

Clean 函数在底层处理时，是这样的，遇到了第一个 `/..` 会回退到前一个非 `/..` 并将其删除，紧接着遍历下一个；如果再次遇到 `/..` 做相同的处理。

# 3 Dir

`func Dir(path string) string`：函数会返回 path 最后一个元素之前的所有部分，通常情况下是 path 路径的目录。**首先会使用 Split 删除掉最后的元素，紧接着使用 Clean 函数，最后删除末尾的斜线**。

如果入参 path 是空字符串，Dir 返回 `.`；如果类似于 `//////xxx`，Dir 返回单个斜线。其他情况下，返回的 path 路径末尾不会是斜线。

~~~go
func TestDir(t *testing.T) {
	fmt.Println(path.Dir("/a/b/c"))
	fmt.Println(path.Dir("a/b/c"))
	fmt.Println(path.Dir("/a/"))
	fmt.Println(path.Dir("a/"))
    fmt.Println(path.Dir("/////////login/"))
	fmt.Println(path.Dir("/"))
	fmt.Println(path.Dir(""))
}

/a/b
a/b
/a
a
/login
/
.
~~~

# 4 Join

`func Join(elem ...string) string`：连接入参的 []string 值元素，组合成单一的 Path 路径，其中元素之间使用斜线分隔开。若元素为空，会忽略该元素。

Join 函数的结果会**经过 Clean 函数**处理，符合其规则。如果入参为空，或者所有元素都为空字符字符串，其返回值也是空字符串。

~~~go
func TestJoin(t *testing.T) {
	fmt.Println(path.Join("a", "b", "c"))
	fmt.Println(path.Join("a", "b/c"))
	fmt.Println(path.Join("a/b", "c"))

	fmt.Println(path.Join("a/b", "../../../xyz"))

	fmt.Println(path.Join("", ""))
	fmt.Println(path.Join("a", ""))
	fmt.Println(path.Join("", "a"))
}

a/b/c
a/b/c
a/b/c
../xyz

a
a
~~~

# 5 Ext



# 6 IsAbs



# 7 Match



# 8 Split





# 9 源代码使用示例

在 Gin 框架中，使用 GET 添加路由：

~~~go
// GET is a shortcut for router.Handle("GET", path, handle).
func (group *RouterGroup) GET(relativePath string, handlers ...HandlerFunc) IRoutes {
    // relativePath 入参，比如类似于 /login
	return group.handle(http.MethodGet, relativePath, handlers)
}

func (group *RouterGroup) handle(httpMethod, relativePath string, handlers HandlersChain) IRoutes {
	absolutePath := group.calculateAbsolutePath(relativePath)
	handlers = group.combineHandlers(handlers)
	debugPrint("absolutePath:%s", absolutePath)

	group.engine.addRoute(httpMethod, absolutePath, handlers)
	return group.returnObj()
}

func (group *RouterGroup) calculateAbsolutePath(relativePath string) string {
    // 默认情况下 group.basePath 为 / 
	debugPrint("group.basePath:%s, relativePath:%s", group.basePath, relativePath)
	return joinPaths(group.basePath, relativePath)
}

func joinPaths(absolutePath, relativePath string) string {
	if relativePath == "" {
		return absolutePath
	}
	// 调用 path 标准库下的 Join 函数，将 absolutePath 和 relativePath 连接起来，组成一个 path 路径
	finalPath := path.Join(absolutePath, relativePath)
    // 经过 path.Join 处理后，path 末尾没有 /
	if lastChar(relativePath) == '/' && lastChar(finalPath) != '/' {
		return finalPath + "/"
	}
	return finalPath
}
~~~

