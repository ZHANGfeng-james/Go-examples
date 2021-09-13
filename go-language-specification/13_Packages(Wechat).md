内容主要解答这样的疑惑：

1. .go 源文件所在的目录名和其中的 package 名有什么关系？
2. .go 源文件中的 PackageClause 和 import 声明怎么使用，具体字段有什么含义？

Go 程序（一个可运行的应用程序，或者是一个功能性的 Module）是由一系列链接在一起的 **Package** 构成的（意思就是说：Go 程序是通过 package 组织起来的）。而 Package 则是由**一个或多个 .go 源代码文件**组成，在这些 .go 源代码文件中可以声明属于该 Package 的常量、类型、变量、函数，并且声明的这些标识符是可以在同一个包下相互能访问到的。这些声明的元素可以被声明为可导出的，这样就可以在其他 Package 中访问。

# 1 Source file organization

**源文件 .go 的组织形式**

每一个 .go 源文件都会在其开头的部分（，使用 package clause）标识出当前文件所属的 Package 名。紧跟这个 package clause 后的是 import 声明，表示会在当前 .go 源文件中会使用到的外部包；紧跟这个 import 声明的是一系列的声明，包括函数、类型、变量和常量。

~~~go
SourceFile       = PackageClause ";" { ImportDecl ";" } { TopLevelDecl ";" } .
~~~

上面的这个 `;` 部分实际上是可以被 Syntax 语义解析，是可以省略的。

# 2 Package clause

包子句

一个包子句写在每一个 .go 源文件的开头位置，用来指出当前 .go 文件所属的 package：

~~~go
PackageClause  = "package" PackageName .
PackageName    = identifier .
~~~

.go 源文件的 package clause 中 PackageName 部分一定不能是空白标识符：`_`，比如

~~~go
package math
~~~

一系列使用相同 PackageName 的 .go 文件构成了整个 Package 的实现。但有一个约束在于：实现**同一个 Package** 的所有 .go 源文件必须放在**同一个目录**下，反过来讲，同一个目录下的所有 .go 源文件中，其 PackageName 必须是一样的。 

# 3 Import declarations

导入声明

一个导入声明包含的 2 个元素分别是：PackageName 和 ImportPath。

导入声明说明了当前 .go 源文件包含了被导入 package 中的声明标识符（比如声明的常量、函数等），允许访问该包下的导出标识符。在导入声明中，命名了一个标识符（PackageName）用于访问包中的内容，命名了一个 ImportPath 表示导入的包。

~~~go
ImportDecl       = "import" ( ImportSpec | "(" { ImportSpec ";" } ")" ) .
ImportSpec       = [ "." | PackageName ] ImportPath .
ImportPath       = string_lit .
~~~

此处的 PackageName 被使用作为限定标识符，用以访问**被导入包**下的**导出标识符**。限定标识符的组成：

~~~go
QualifiedIdent = PackageName "." identifier .
~~~

这些导出标识符应该是被声明在**文件块**中。如果省略 PackageName，则 PackageName 的默认值就是导入包的 package 子句中指定的标识符（也就是下述 PackageClause 中的 PackageName）。

~~~go
PackageClause  = "package" PackageName .
PackageName    = identifier .
~~~

如果把 PackageName 显式地写为 `.`，则就可以在当前源文件中直接访问被导入包的导出标识符，而不需要显式使用包限定符。

ImportPath 的值和当前导入包的实现相关，其值通常是完整文件名的字符串，并且可能和已安装软件包的存储仓库相关。

假设我们已经编译了一个包含有 `package math` 的 Package，该包的功能导出了 Sin 函数，同时在 `lib/math` 路径下安装了这个已编译的包。下面这些方式，可以用于访问 Sin 函数：

~~~go
Import declaration          Local name of Sin

import   "lib/math"         math.Sin
import m "lib/math"         m.Sin
import . "lib/math"         Sin
~~~

一个导入声明表明了导入和导入包之间的依赖关系。包直接或间接导入自身（包括递归导入包的情况），或直接导入包但不引用其任何导出标识符，这些都是非法的。要仅处于副作用（初始化）导入软件包，可以使用空白标识符作为显式的 PackageName。比如：

~~~go
import _ "lib/math"
~~~

# 4 使用场景

在实际工程中，有 2 种情况：

1. PackageName 和目录名不一致；
2. 不同目录下定义了相同的 PackageName。

使用这样的工程结构做测试：

~~~go
G:\Go\SyntaxTest>tree /F
G:.
│  go.mod
│  main.go
│
├─selector
│      init.go
│      init_test.go
│
└─slice
        init.go
        other.go
~~~

对于**第一种情况**：PackageName 和目录名不一致，也就是说在目录 slice 下创建的 init.go 和 other.go 中的 Package clause 的 PackageName 和目录名 slice 不相同，比如使用的 PackageName 是 reslice。在导包时，可以这样使用：

~~~go
package main

import reslice "opensource.com/syntaxtest/slice"

func main() {
	reslice.InitTest5()
}
~~~

当然 init.go 和 other.go 的 PackageName 必须是一样的。而且 `import reslice "opensource.com/syntaxtest/slice"` 中 PackageName，也可以和 slice 目录下 init.go 中 PackageClause 的 PackageName 不一样！

对于**第二种情况**：不同目录下定义了相同的 PackageName。比如在 selector 目录下 init.go 文件和 slice 目录下的 init.go 文件，都使用相同的 PackageClause：`package reslice`。可以像下面这种情况使用：

~~~go
package main

import (
	resliceB "opensource.com/syntaxtest/selector"
	resliceA "opensource.com/syntaxtest/slice"
)

func main() {
	resliceA.InitTest5()
	resliceB.Test()
}
~~~

分别使用 resliceA 和 resliceB 区分了不同的 package。