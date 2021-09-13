**Tutorial 系列文章**介绍 Go 的一些**主要功能**：

1. 创建模块（module）：可供其他模块调用的功能代码集合
2. 调用模块功能；
3. 错误处理；
4. Go 的动态数组：Slice 切片；
5. Map 中的 key-value 对；
6. 单元测试；
7. 编译和应用程序安装。

# 1 创建模块

Go 中的 Module 定义的是一个或多个相关程序包（Package），包含了一组离散且具有一定用途的函数。

Go 代码被分组到包（Package）中，而软件包（Package）则又被分组到模块（Module）中。包（Package）的模块（Module）指定了 Go 运行代码所需要的**上下文**，包括编写代码的 Go 版本及其所需的其他模块集合。

当源模块有改动时，源模块作者会发布新版本。模块的调用者会导入模块的更新文件，并在正式导入到产品前测试该新模块。

比如下述示例，创建了 greetings 模块（对应了平台下的目录），并在该模块下撰写代码：

~~~go
C:\Users\Developer\hello>cd ..

C:\Users\Developer>mkdir greetings

C:\Users\Developer>cd greetings

C:\Users\Developer\greetings>go mod init example.com/greetings
go: creating new go.mod: module example.com/greetings
~~~

go mod init 指令时（创建了 go.mod 文件），给定了参数 example.com/greetings 作为**模块路径**。在生产代码中，可以使用该路径下载对应模块。go.mod 文件**将当前代码标识为**可以在其他代码中使用的**模块**，包含了模块名以及支持的 Go 版本信息。另外，如果导入了第三方模块，go.mod 文件会列出第三方模块的指定版本信息。go.mod 让软件的构建具备了可复制性，并可直接控制要使用的模块版本。

在 greetings 模块中创建 greentings.go 文件并增加如下代码：

~~~go
package greetings

import "fmt"

// Hello returns a greeting for the named person.
func Hello(name string) string {
    // Return a greeting that embeds the name in a message.
    message := fmt.Sprintf("Hi, %v. Welcome!", name)
    return message
}
~~~

上述代码中声明了 greetings 包（Package），用于**封装相关的功能代码**。

![](./pics/Snipaste_2020-12-02_09-33-53.png)

在 Go 中，:= 运算符是一种用于在一行代码中**声明和初始化变量的快捷方式**（Go 使用右侧的值来确定变量的类型）。或者也可以扩展开来：

~~~go
var message string
message = fmt.Sprintf("Hi, %v. Welcom!", name)
~~~

# 2 调用模块

复用之前的 hello 模块，并**引入 greetings 模块的功能代码**：

~~~go
package main

import(
    "fmt"
    "rsc.io/quote"
    "example.com/greetings"
)

func main(){
    fmt.Println(quote.Go())
    // Get a greeting message and print it.
    message := greetings.Hello("Gladys")
	fmt.Println(message)
}
~~~

上述代码中，**声明了 main 包**，在 Go 中，作为应用程序执行的代码**必须放置到 main 包下**。同时**导入了 3 个包**："fmt"、"rsc.io/quote"、"example.com/greetings"，使用这种方法可以访问到这些包下定义的功能方法。

对于当前的 hello 目录而言，若没有 go.mod 文件，则可以认为并不是一个 Module，而仅仅只是一个目录而已。因此为了让 Go 能将 hello 识别为 Module，需要执行 go mod init hello 并创建了 go.mod 文件（删除原先 hello 目录下的 go.mod 和 go.sum 文件）：

~~~go
C:\Users\Developer\hello>go mod init hello
go: creating new go.mod: module hello
~~~

那接下来的逻辑就是要让 hello 模块在运行时能够找到 example.com/greetings 模块！

**对于生产用途，可以将模块发布在公司内部或者网络服务器上，然后通过 Go 命令下载**。对于当前情况（模块功能存在于本地文件系统中），需要调整这种方式以便让 hello 模块能访问本地的 example.com/greetings。可这样修改：

~~~go
module hello

go 1.15

replace example.com/greetings => ../greetings
~~~

replace 指令告诉 Go 使用 `../greetings` 替换模块地址 `example.com/greetings`。

运行 go build 指令，让 Go 去定位模块并将作为依赖项添加到 go.mod 文件中

~~~go
C:\Users\Developer\hello>go build
go: finding module for package rsc.io/quote
go: found example.com/greetings in example.com/greetings v0.0.0-00010101000000-000000000000
go: found rsc.io/quote in rsc.io/quote v1.5.2
~~~

运行后 go.mod 文件内容变更为：

~~~go
module hello

go 1.15

replace example.com/greetings => ../greetings

require (
	example.com/greetings v0.0.0-00010101000000-000000000000
	rsc.io/quote v1.5.2
)
~~~

在 build 指令执行过程中，Go 会定位到 `../greetings` 目录，并增加 require 指令用以表明 hello 是依赖于 `example.com/greetings` 的。当 hello.go 中导入 greetings 包时，便创建了该依赖关系。replace 指令告诉 Go 在哪里可以找到 greetings 模块。对于已发布的 Module，在 go.mod 文件中需要删除 replace 指令，直接使用末尾带有版本信息的 require 指令代替。`require example.com/greetings v1.1.0`。

运行了 go build 指令后，编译出了 hello.exe 可执行文件。

~~~go
C:\Users\Developer\hello>hello.exe
Don't communicate by sharing memory, share memory by communicating.
Hi, Gladys. Welcome!
~~~

# 3 错误处理

==处理错误是可靠代码的基本特征！==

修改 greetings.go 代码：

~~~go
package greetings

import (
    "errors"
    "fmt"
)

// Hello returns a greeting for the named person.
func Hello(name string) (string, error) {
    // If no name was given, return an error with a message.
    if name == "" {
        return "", errors.New("empty name")
    }

    // If a name was received, return a value that embeds the name 
    // in a greeting message.
    message := fmt.Sprintf("Hi, %v. Welcome!", name)
    return message, nil
}
~~~

Hello 函数会返回 2 个类型值：String 和 error。调用者在调用后，会检查函数的返回结果，用以判断结果是否出错。在参数为空的情况时，使用 errors.New() 创建一个 error 对象，nil 则表示没有错误。

修改 hello.go 代码：

~~~go
package main

import(
    "fmt"
    "log"
    "rsc.io/quote"
    "example.com/greetings"
)

func main(){
    fmt.Println(quote.Go())

    // Set properties of the predefined Logger, including
    // the log entry prefix and a flag to disable printing
    // the time, source file, and line number.
    log.SetPrefix("greetings: ")
    log.SetFlags(0)

    // Request a greeting message.
    message, err := greetings.Hello("")
    // If an error was returned, print it to the console and
    // exit the program.
    if err != nil {
        log.Fatal(err)
    }

    // If no error was returned, print the returned message
    // to the console.
    fmt.Println(message)
}
~~~

运行上述程序：

~~~go
C:\Users\Developer\hello>go build

C:\Users\Developer\hello>hello.exe
Don't communicate by sharing memory, share memory by communicating.
greetings: empty name

C:\Users\Developer\hello>go run hello.go
Don't communicate by sharing memory, share memory by communicating.
greetings: empty name
exit status 1
~~~

实际上，这就是 Go 中错误处理的工作方式：**将错误作为值返回，以便调用者可以检查是否出错**。

# 4 切片

Slice 类似于数组，不同之处在于：在添加和删除元素时，会动态调整大小。

修改 greetings 代码如下：

~~~go
package greetings

import (
    "errors"
    "fmt"
    "math/rand"
    "time"
)

// Hello returns a greeting for the named person.
func Hello(name string) (string, error) {
    // If no name was given, return an error with a message.
    if name == "" {
        return "", errors.New("empty name")
    }

    // If a name was received, return a value that embeds the name 
    // in a greeting message.
    message := fmt.Sprintf(randomFormat(), name)
    return message, nil
}

func init(){
    rand.Seed(time.Now().UnixNano())
}

// randomFormat returns one of a set of greeting messages. The returned
// message is selected at random.
func randomFormat() string {
    // A slice of message formats.
    formats := []string{
        "Hi, %v. Welcome!",
        "Great to see you, %v!",
        "Hail, %v! Well met!",
    }

    // Return a randomly selected message format by specifying
    // a random index for the slice of formats.
    return formats[rand.Intn(len(formats))]
}
~~~

在上述代码中引入了 randomFormat()，用于返回随机的格式字符串。该方法是以小写字符开头的，意味着其作用域是本 Package 中，在包外是无法访问到的。

在 randomFormat() 中声明了 formats 切片对象，可以使用 []string 方式**省略其长度**，表示的含义：告诉 Go 当前定义的是**能够动态调整数组大小的切片对象**。

在执行 greetings 包下的代码时，Go 会在初始化全局变量之后，自动调用 init()。

修改调用者代码：greetings.Hello(“Gladys”) 多次运行程序可得到如下结果：

~~~go
C:\Users\Developer\hello>go run hello.go
Don't communicate by sharing memory, share memory by communicating.
Great to see you, Gladys!

C:\Users\Developer\hello>go run hello.go
Don't communicate by sharing memory, share memory by communicating.
Hail, Gladys! Well met!

C:\Users\Developer\hello>go run hello.go
Don't communicate by sharing memory, share memory by communicating.
Hi, Gladys. Welcome!
~~~

# 5 Map

上节的程序实现了调用 Hello(name string) 返回单一字符串的功能，本节则是实现多输入多输出的功能。

为了实现上述功能，需要在调用 Hello 时传递一个集合对象。因此，需要修改 greetings 下的 Hello 方法的签名。如果 greetings 模块已经发布，而且调用者已经使用了 Hello 方法构建的应用程序，此次修改必然回调会导致程序崩溃。最好的修改方法是：==给新功能启用一个新名字==，这是**一种能保持软件向下兼容的想法**。

> 在计算机中指在一个程序或者类库更新到较新的版本后，用旧的版本程序创建的文档或系统仍能被正常操作或使用，或在旧版本的类库的基础上开发的程序仍能正常编译运行的情况。

修改 greetings 代码如下：

~~~go
package greetings

import (
    "errors"
    "fmt"
    "math/rand"
    "time"
)

// Hello returns a greeting for the named person.
func Hello(name string) (string, error) {
    // If no name was given, return an error with a message.
    if name == "" {
        return name, errors.New("empty name")
    }
    // Create a message using a random format.
    message := fmt.Sprintf(randomFormat(), name)
    return message, nil
}

// Hellos returns a map that associates each of the named people
// with a greeting message.
func Hellos(names []string) (map[string]string, error) {
    // A map to associate names with messages.
    messages := make(map[string]string)
    // Loop through the received slice of names, calling
    // the Hello function to get a message for each name.
    for _, name := range names {
        message, err := Hello(name)
        if err != nil {
            return nil, err
        }
        // In the map, associate the retrieved message with 
        // the name.
        messages[name] = message
    }
    return messages, nil
}

// Init sets initial values for variables used in the function.
func init() {
    rand.Seed(time.Now().UnixNano())
}

// randomFormat returns one of a set of greeting messages. The returned
// message is selected at random.
func randomFormat() string {
    // A slice of message formats.
    formats := []string{
        "Hi, %v. Welcome!",
        "Great to see you, %v!",
        "Hail, %v! Well met!",
    }

    // Return one of the message formats selected at random.
    return formats[rand.Intn(len(formats))]
}
~~~

上述程序中新增了 Hellos 方法，其参数不再是单一的名字，而是包含很多名字的切片；相应的，其方法返回值也从单一的字符串变为了 Map。在 Go 中，make(map[key-type]value-type) 初始化 Map 对象。

~~~go
for _, name := range names {
    // 调用 Hello 函数，其结果反馈 message 和 err 对象
    message, err := Hello(name)
    if err != nil {
        // 若出错，则将 messages 置位 nil
        return nil, err
    }
    // In the map, associate the retrieved message with the name.
    // key - value 键值对赋值
    messages[name] = message
}
~~~

在 for 循环中，range 返回了 2 个值：循环中当前元素的索引和该元素的值（拷贝）。如果不需要位置索引，可以使用下划线忽略。

修改 hello.go 程序代码：

~~~go
package main

import(
    "fmt"
    "log"
    "rsc.io/quote"
    "example.com/greetings"
)

func main(){
    fmt.Println(quote.Go())

    // Set properties of the predefined Logger, including
    // the log entry prefix and a flag to disable printing
    // the time, source file, and line number.
    log.SetPrefix("greetings: ")
    log.SetFlags(0)

    names :=[]string{"Gladys", "Samantha", "Darrin"}

    // Request a greeting message.
    messages, err := greetings.Hellos(names)
    // If an error was returned, print it to the console and
    // exit the program.
    if err != nil {
        log.Fatal(err)
    }

    // If no error was returned, print the returned message
    // to the console.
    fmt.Println(messages)
}
~~~

程序运行结果：

~~~go
C:\Users\Developer\hello>go build

C:\Users\Developer\hello>hello.exe
Don't communicate by sharing memory, share memory by communicating.
map[Darrin:Hi, Darrin. Welcome! Gladys:Hail, Gladys! Well met! Samantha:Hi, Samantha. Welcome!]

C:\Users\Developer\hello>hello.exe
Don't communicate by sharing memory, share memory by communicating.
map[Darrin:Great to see you, Darrin! Gladys:Hail, Gladys! Well met! Samantha:Hi, Samantha. Welcome!]

C:\Users\Developer\hello>hello.exe
Don't communicate by sharing memory, share memory by communicating.
map[Darrin:Hail, Darrin! Well met! Gladys:Great to see you, Gladys! Samantha:Great to see you, Samantha!]
~~~

# 6 单元测试

在开发工程中，测试代码会暴露出修改过程中出现的 bug。

单元测试是 Go 的内置支持功能，具体来说，是使用命名约定、Go testing 包以及 go test 指令，可以快速编写和执行测试程序。

在 greetings 目录下创建 greetings_test.go 文件，**以 _test.go 结尾的文件名**告诉 go test 指令该文件包含单元测试方法（具备单元测试功能）。

~~~go
package greetings

import (
    "testing"
    "regexp"
)

// TestHelloName calls greetings.Hello with a name, checking 
// for a valid return value.
func TestHelloName(t *testing.T) {
    name := "Gladys"
    want := regexp.MustCompile(`\b`+name+`\b`)
    msg, err := Hello("Gladys")
    if !want.MatchString(msg) || err != nil {
        t.Fatalf(`Hello("Gladys") = %q, %v, want match for %#q, nil`, msg, err, want)
    }
}

// TestHelloEmpty calls greetings.Hello with an empty string, 
// checking for an error.
func TestHelloEmpty(t *testing.T) {
    msg, err := Hello("")
    if msg != "" || err == nil {
        t.Fatalf(`Hello("") = %q, %v, want "", error`, msg, err)
    }
}
~~~

在 greetings_test.go 文件中，创建了 2 个函数 TestHelloName 和 TestHelloEmpty 用来测试 Hello 功能。

测试函数的命名使用 TestName 格式，其中 Name 是特定于测试的，同时使用了指向 testing 包下的 testing.T 指针作为参数，可以使用该参数的方法获取测试报告和日志。

在 greetings 目录下运行 go test 指令，会执行 xxx_test.go 文件中的 Testxxx 函数：

~~~go
C:\Users\Developer\greetings>go test
PASS
ok      example.com/greetings   0.251s

C:\Users\Developer\greetings>go test -v
=== RUN   TestHelloName
--- PASS: TestHelloName (0.00s)
=== RUN   TestHelloEmpty
--- PASS: TestHelloEmpty (0.00s)
PASS
ok      example.com/greetings   0.270s
~~~

那现在模拟单元测试失败的案例，修改 greeting.go 源代码：

~~~go
// Hello returns a greeting for the named person.
func Hello(name string) (string, error) {
    // If no name was given, return an error with a message.
    if name == "" {
        return name, errors.New("empty name")
    }
    // Create a message using a random format.
    // message := fmt.Sprintf(randomFormat(), name)
    message := fmt.Sprint(randomFormat())
    return message, nil
}
~~~

这次，仅运行 go test，该结果仅会输出测试失败的结果，这在执行大量单元测试时非常有效。

~~~go
C:\Users\Developer\greetings>go test -v
=== RUN   TestHelloName
    greetings_test.go:15: Hello("Gladys") = "Great to see you, %v!", <nil>, want match for `\bGladys\b`, nil
--- FAIL: TestHelloName (0.00s)
=== RUN   TestHelloEmpty
--- PASS: TestHelloEmpty (0.00s)
FAIL
exit status 1
FAIL    example.com/greetings   0.256s

C:\Users\Developer\greetings>go test
--- FAIL: TestHelloName (0.00s)
    greetings_test.go:15: Hello("Gladys") = "Hail, %v! Well met!", <nil>, want match for `\bGladys\b`, nil
FAIL
exit status 1
FAIL    example.com/greetings   0.267s
~~~

# 7 编译和安装程序

切换当前目录到 hello 下，执行 go list -f '{{.Target}}' 输出如下内容：

~~~
C:\Users\Developer\hello>go list -f '{{.Target}}'
'C:\Users\Developer\go\bin\hello.exe'
~~~

上述结果表明，go install 指令会将 hello.exe 安装到 `C:\Users\Developer\go\bin` 路径下。

将上述安装路径添加到 path 路径下：

~~~go
C:\Users\Developer\hello>set PATH=%PATH%;C:\Users\Developer\go\bin

C:\Users\Developer\hello>path
PATH=C:\Program Files (x86)\NetSarang\Xftp 6\;C:\WINDOWS\system32;C:\WINDOWS;C:\WINDOWS\System32\Wbem;C:\WINDOWS\System32\WindowsPowerShell\v1.0\;C:\WINDOWS\System32\OpenSSH\;C:\Program Files\Microsoft SQL Server\110\Tools\Binn\;E:\Java\jdk1.8.0_144\bin;E:\Java\jdk1.8.0_144\jre\bin;e:\Git\cmd;E:\Android\StudioSDK\platform-tools;E:\Android\StudioSDK\tools;E:\Android\StudioSDK\build-tools\28.0.2;E:\Program Files\nodejs\;E:\Program Files\wkhtmltopdf\bin;C:\Program Files (x86)\Google\Chrome\Application;C:\Users\Developer\AppData\Local\Programs\Python\Python37\Scripts\phantomjs-2.1.1-windows\bin;E:\TortoiseGit\bin;C:\Program Files (x86)\Pandoc\;E:\Android\StudioSDK\ndk-bundle;C:\Go\bin;C:\Users\Developer\AppData\Local\Microsoft\WindowsApps;C:\Users\Developer\AppData\Roaming\npm;e:\JetBrains\PyCharm2018.3.3\bin;E:\apache-maven-3.6.1\bin;E:\Program Files\mingw-w64\x86_64-8.1.0-win32-seh-rt_v6-rev0\mingw64\bin;E:\SQLite;E:\apache-maven-3.6.3\bin;C:\Program Files (x86)\WinRAR;C:\Users\Developer\go\bin;C:\Users\Developer\go\bin
~~~

这样，在 go install 后，不管当前执行的目录在何处，都可以通过 path 路径找到可执行程序：

~~~go
C:\Users\Developer\hello>go install

C:\Users\Developer\hello>hello
Don't communicate by sharing memory, share memory by communicating.
map[Darrin:Great to see you, Darrin! Gladys:Hail, Gladys! Well met! Samantha:Great to see you, Samantha!]

C:\Users\Developer\hello>cd ..

C:\Users\Developer>hello
Don't communicate by sharing memory, share memory by communicating.
map[Darrin:Hail, Darrin! Well met! Gladys:Hail, Gladys! Well met! Samantha:Hail, Samantha! Well met!]
~~~

![](./pics/Snipaste_2020-12-02_16-19-47.png)

或者是使用 go env 指令改变 GOBIN 变量值，修改安装路径：`go env -w GOBIN=C:\Users\Developer\go\bin`
