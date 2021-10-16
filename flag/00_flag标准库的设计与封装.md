



从**结构体定义**开始：

~~~go
// A Flag represents the state of a flag.
type Flag struct {
	Name     string // name as it appears on command line
	Usage    string // help message
	Value    Value  // value as set
	DefValue string // default value (as text); for usage message
}
~~~

Flag 结构体为什么定义成这个样子？



~~~go
// A FlagSet represents a set of defined flags. The zero value of a FlagSet
// has no name and has ContinueOnError error handling.
//
// Flag names must be unique within a FlagSet. An attempt to define a flag whose
// name is already in use will cause a panic.
type FlagSet struct {
	// Usage is the function called when an error occurs while parsing flags.
	// The field is a function (not a method) that may be changed to point to
	// a custom error handler. What happens after Usage is called depends
	// on the ErrorHandling setting; for the command line, this defaults
	// to ExitOnError, which exits the program after calling Usage.
	Usage func()

	name          string
	parsed        bool
	actual        map[string]*Flag
	formal        map[string]*Flag
	args          []string // arguments after flags
	errorHandling ErrorHandling
	output        io.Writer // nil means stderr; use Output() accessor
}
~~~

什么是 FlagSet？有什么含义？

命令格式 `go <command> [arguments]` 中的 command 就是**子命令**，因此 bug 和 build 就是**不同功能的子命令**。go 程序在运行时（**go 相当于是一个应用程序**），会首先**识别出对应的子命令**，紧接着**再去解析子命令之后的参数**。

FlagSet 中的 name 字段对应的就是 command 子命令。