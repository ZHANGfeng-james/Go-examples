Go 语言中的 block，其含义在于：语义的范围，可简称为“代码块”。

一个代码块，在形式上是使用“{}”包围起来的，其中包含有一些 declarations（声明） 和 statements（语句）。其格式可为：

~~~go
Block = "{" StatementList "}" .
StatementList = { Statement ";" } .
~~~

在 Go 代码中，除了显式地写明“{}”表示代码块，还有隐藏的代码块：

* Go 源代码在一个**全局代码块**中；

* 每一个 package 中的源代码都处在一个**包级代码块**中；

* 每一个源代码文件都处在一个**源文件级代码块**中；

* if、for、switch **语句**有其自身的**隐式代码块**；

  此处的 if、for、switch 语句**自身的隐式代码块**主要表现在 `StatementList = { Statement ";" } .` 中，也就是说：for、if、switch 关键字后面可接上 `;` 这样的语句。同时，可以在这些语句中声明变量，这些声明的变量可以“覆盖” main 代码块中的变量。

* 在 switch 和 select 语句中的**子句**都有其**隐式代码块**。

  此处的隐式代码块，是指 case 关键字之后的内容。

# for

for 语句有其自身的隐式代码块：

~~~go
var index int8
index = 11
fmt.Println("index:", index)
// for
for index := 0; index < 10; index++ {
	fmt.Println("for loop inner! index:", index)
	break
}

PS G:\Go\go_developer_roadmap\OpenSource\LoadGenerator> go run main.go
index: 11
for loop inner! index: 0
~~~

也就是说在 for 隐式代码块中，index 实际上指代的是 `index := 0` 声明的变量，而不是在 main 作用域中声明的变量。或者可以这么说，for 隐式代码块中的 index 变量覆盖了 main 中 index 变量。

# if

if 语句有其自身的隐式代码块：

~~~go
var flag bool
flag = false
fmt.Println("flag:", flag)
// if
if flag := getValue(); flag {
	fmt.Println("if statement inner, flag:", flag)
}

func getValue() bool {
	return true
}

PS G:\Go\go_developer_roadmap\OpenSource\LoadGenerator> go run main.go
flag: false
if statement inner, flag: true
~~~

同理，if 也是有其自身的隐式代码块的，if 隐式代码块中定义的 `flag` “覆盖”了 main 代码块中定义的 flag 变量。

# switch

~~~go
package main

import "fmt"

func main() {

	var value int8
	value = 0
	fmt.Println("value:", value)
	// switch
	switch value := getValue(); value {
	case 0:
		fmt.Println("switch statement inner, value:", value)
	case 1:
		fmt.Println("switch statement inner, value:", value)
		fmt.Printf("%T.\n", value)

		value := true
		fmt.Println(value)
	}

}

func getValue() int {
	return 1
}
PS G:\Go\go_developer_roadmap\OpenSource\LoadGenerator> go run main.go
value: 0
switch statement inner, value: 1
int.
true
~~~

同理，此处的 if、for、switch 语句**自身的隐式代码块**主要表现在 `StatementList = { Statement ";" } .` 中，也就是说：for、if、switch 关键字后面可接上 `;` 这样的语句。同时，可以在这些语句中声明变量，这些声明的变量可以“覆盖” main 代码块中的变量。

另外，从 case 1 的子句（Clause）中，可以再次覆盖 switch statement 中声明的变量，从而再一次形成了一个隐式代码块。

# select

~~~go
package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan struct{})
	go func() {
		ch <- struct{}{}
	}()

	time.Sleep(time.Millisecond * 100)
	value := true
	fmt.Printf("%T, value:%v.\n", value, value)
	select {
	case value, ok := <-ch:
		fmt.Printf("value:%s, ok:%v.\n", value, ok)
	default:
		fmt.Println("default clause inner!")
	}
}
PS G:\Go\go_developer_roadmap\OpenSource\LoadGenerator> go run main.go
bool, value:true.
value:{}, ok:true.
~~~

select 的隐式代码块主要表现在其各个 case clause 中，如上所述，value 变量在 case clause 子句中“覆盖”了 main 中的 value 变量。

