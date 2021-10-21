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

对于一个 Flag 来说，对应的使用场景是：`go run main.go -name="Katyusha"`，在运行指令中 `name` 就是一个 Flag。对应的这个抽象出来的 Flag 实体，有对应的 Name、Usage、Value 等属性。

**再往“前”探索**，程序在执行时，引入 flag 库的过程中，实际上做了如下初始化操作：

~~~go
// CommandLine is the default set of command-line flags, parsed from os.Args.
// The top-level functions such as BoolVar, Arg, and so on are wrappers for the
// methods of CommandLine.
var CommandLine = NewFlagSet(os.Args[0], ExitOnError)

func init() {
	// Override generic FlagSet default Usage with call to global Usage.
	// Note: This is not CommandLine.Usage = Usage,
	// because we want any eventual call to use any updated value of Usage,
	// not the value it has when this line is run.
	CommandLine.Usage = commandLineUsage
}
~~~

其中 `os.Args[0]` 是当前执行进程的名称，比如：`/var/folders/3w/0pprfk1x13lfrt7vdxryd75r0000gn/T/go-build127194955/b001/exe/main`。标准库 flag 会**默认创建一个 CommandLine 这个 FlagSet，接下来看这个结构体**：

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

什么是 FlagSet？有什么含义？从 FlagSet 的定义上来看：**是一个定义的 Flag 的集合**。也就是，在 FlagSet 下面有一系列的 Flag。FlagSet 和 Flag 之间是一种**层级关系**，FlagSet 包含了 Flag。

比如：

~~~go
ant@MacBook-Pro usage % go run main.go --age=32
2021/10/21 10:13:40 main.go:15: /var/folders/3w/0pprfk1x13lfrt7vdxryd75r0000gn/T/go-build933258456/b001/exe/main
2021/10/21 10:13:40 main.go:53: age:32
2021/10/21 10:13:40 main.go:60: GetInt:32
2021/10/21 10:13:40 main.go:64: GetInt error:flag accessed but not defined: port
~~~

就是直接使用了 `/var/folders/3w/0pprfk1x13lfrt7vdxryd75r0000gn/T/go-build933258456/b001/exe/main` 这个已经被创建的 FlagSet。