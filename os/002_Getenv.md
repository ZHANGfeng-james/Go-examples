`func Getenv(key string) string`：解析（获取） key 指定的**环境变量值**。

> 环境变量，一般是指在操作系统中**用来指定操作系统运行环境的一些参数**，如临时文件夹位置和系统文件夹位置等。
>
> 环境变量是在操作系统中**一个具有特定名字的对象（一般是 key-value）**，它包含了**一个或者多个应用程序所使用到的信息**。

函数返回值类型是 string，如果该 key 指定的环境变量不存在（未被设置），其值是空字符串。为了区分这种情况，可以使用 LookupEnv 函数作为**辅助查询**。

~~~go
package main

import (
	"fmt"
	"os"
)

func main() {
	value, ok := os.LookupEnv("GOPATH")
	if ok {
		fmt.Printf("Env Present, value is %s\n", value)
	}

	env := "NAME"
	gopath := os.Getenv(env)
	fmt.Printf("%s value is <%s>\n", env, gopath)
}
PS E:\go_developer_roadmap\ProgrammingLanguage\Go Standard Interface\GoUsage> go run main.go
Env Present, value is C:\Users\Administrator\go
NAME value is <>
~~~

与之相关的是 `func LookupEnv(key string) (string, bool)`：其返回值 bool 的含义明确了当前进程**是否设置了指定的环境变量**。来看看示例代码：

~~~go
package main

import (
	"fmt"
	"os"
)

func main() {
	show := func(key string) {
		val, ok := os.LookupEnv(key)
		if !ok {
			fmt.Printf("%s not set\n", key)
		} else {
			fmt.Printf("%s=%s\n", key, val)
		}
	}
	os.Setenv("SOME_KEY", "value")
	os.Setenv("EMPTY_KEY", "")

	show("SOME_KEY")
	show("EMPTY_KEY")
	show("eMPTY_KEY")
	show("MISSING_KEY")
}
PS E:\go_developer_roadmap\ProgrammingLanguage\Go Standard Interface\GoUsage> go run main.go
SOME_KEY=value
EMPTY_KEY=
eMPTY_KEY=
MISSING_KEY not set
~~~

通过实际测试发现：**Windows** 运行操作系统上，**环境变量是大小写不敏感的**，比如上述 `eMPTY_KEY` 和 `EMPTY_KEY` 被解读成相同的环境变量。

