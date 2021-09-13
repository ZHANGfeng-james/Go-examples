

flag 标准库**源代码路径**：`./Go/src/flag/` 目录，其中包含：

1. `flag.go`：flag 标准库的具体实现；
2. `flag_test.go`、`export_test.go`、`example_test.go`、`example_value_test.go`：flag 标准库测试程序。

# 功能描述

flag 包**实现了命令行参数的解析**，让开发命令行工具更简单了。何为“命令行参数”？比如我们运行程序时给定指令：`go run main.go -name ant` 其中的 `-name ant` 就是命令行参数。命令行参数的所有内容，可通过 `os.Args` 获取到。

可将 flag 包下封装的内容看作（理解成）是 Go 实现的解析模块，该模块实现的功能就是去解析 `-name ant`、`-name=ant` 等形式的命令行参数值。即根据 `name` 标记的命令行参数，解析得到值 `ant`。同时，Go 将实现解析功能的模块开放给 App 使用，实现了命令行参数的获取、使用。

从接口文档来看，flag 包的实现==分为 2 个层次==：

1. 单一 flag 的解析，比如 `int`、`int64`、`uint`、`uint64`、`float64`、`bool`、`duration`、`string` 类型的解析；
2. `FlagSet` 的解析，实现子命令。同样也实现了上述 8 种类型 flag 值的解析。

另外从 API 上也可以很清晰看出来，有一类是 `*FlagSet` 作为接收者的方法。

| flag 参数 |                            有效值                            |
| :-------: | :----------------------------------------------------------: |
|  字符串   |                          合法字符串                          |
|   整数    |            1234、0664、0x1234等类型，也可以是负数            |
|   浮点    |                          合法浮点值                          |
|   bool    |   1, 0, t, f, T, F, true, false, TRUE, FALSE, True, False    |
|  时间段   | 任何合法的时间段字符串。如”300ms”、”-1.5h”、”2h45m”。合法的单位有”ns”、”us” /“µs”、”ms”、”s”、”m”、”h” |

# 用法及示例程序

~~~go
var nFlag = flag.Int("n", 1234, "help message for flag n")
~~~

先来看几种相似类型的方法，使用 `flag.Int()`\`String()`\`Bool()` 等函数定义 flag，比如上述就定义了一个 `int` 类型的 flag，其名称为 `-n`，值存储在变量 `nFlag` 中，类型为 `*int`。通过解析，可直接 `*nFlag` 获取解析的值。

或者可使用 `Var()` 系列函数将 flag 绑定到变量中：

~~~go
var flagvar int
flag.IntVar(&flagvar, "flagname", 1234, "help message for flagname")
~~~

或者可以**自定义一个用于 flag 的类型**（需满足 Value 接口），并将该类型用于 flag 解析：

~~~go
flag.Var(&flagvar, "name", "help message for flagname")
~~~

对这种 flag，默认值就是该变量的初始值。

在所有 flag 都定义好后，调用 `flag.Parse()` 实现**命令行参数写入注册的 flag 的解析动作**。解析之后，flag 的值可以直接使用。如果使用的是 flag 本身，则其类型就是指针类型；如果绑定在变量上，则可通过变量获取解析值。

~~~go
fmt.Println("ip has value ", *ip)
fmt.Println("flagvar has value ", flagvar)
~~~

解析后，flag 后面的参数可以从 `flag.Args()` 中获取或者用 `flag.Arg(i)` 单独获取。

==命令行 flag 的语法遵循下述规则==：

~~~go
-flag
-flag=x
-flag x // 非 boolean 命令行 flag
~~~

可以使用1个或2个`-`符号，效果是一样的。最后一种格式不能用于 bool 类型的 flag，因为如果有文件名为 0、false 等时，如下命令 `cmd -x *` 其含义会改变。必须使用 `-flag=false` 格式来关闭一个 bool 类型 flag。

flag 解析在第一个非 flag 参数（单个 `-` 不是 flag 参数）之前停止，或者在终止符 `--` 之后停止。

整数 flag 可接收的值是 1234、0664、0x1234，也可以是负数。Boolean flag 可以接收：`1, 0, t, f, T, F, true, false, TRUE, FALSE, True, False` 这样的值。

时间段 flag 接收任何合法的可提供给 `time.ParseDuration` 的输入。

默认的命令行 flag 集被**包水平**的函数控制，`FlagSet` 类型允许程序员定义独立的 flag 集合，例如实现命令行界面下的子命令。

关于 flag 包的使用可参考如下示例代码：

~~~go
package main

import (
	"errors"
	"flag"
	"fmt"
	"strings"
	"time"
)

var gopherType string

func init() {
	const (
		defaultGopher = "pocket"
		usage         = "the variety of gopher"
	)
	flag.StringVar(&gopherType, "gopher_type", defaultGopher, usage)
	// shorthand usage, use the same default value!
	flag.StringVar(&gopherType, "g", defaultGopher, usage+" (shorthand)")

	flag.Var(&intervalFlag, "deltaT", "comma-separated list of intervals to use between events")
}

type interval []time.Duration

var intervalFlag interval

func (i *interval) String() string {
	return fmt.Sprint(*i)
}

func (i *interval) Set(value string) error {
	fmt.Println("input params:" + value)

	// go run main.go -species ant -gopher_type a -g b -deltaT 10s -deltaT 12s 注释下述 if 语句
	if len(*i) > 0 {
		return errors.New("interval flag already set")
	}
	for _, dt := range strings.Split(value, ",") {
		duration, err := time.ParseDuration(dt)
		if err != nil {
			return err
		}
		*i = append(*i, duration)
	}
	return nil
}

func main() {
	var species = flag.String("species", "gopher", "the species we are studying")
	flag.Parse()

	fmt.Printf("%s\n", *species)
	fmt.Printf("%s\n", gopherType)

	fmt.Printf("%s\n", intervalFlag)
}
~~~



# 源代码解析

整个源代码解析的过程可以**分为 2 个层次**：flag、`FlagSet`；整个实现的过程可以**分为 2 个阶段**：注册、解析。

先全局看看定义的类型、接口，定义了动态值（用于存放多种不同类型值）的接口：

~~~go
type Value interface {
    String() string
    Set(string) error
}
~~~

Value 接口的直接使用是在定义结构体 Flag 中：

~~~go
// A Flag represents the state of a flag.
type Flag struct {
	Name     string // name as it appears on command line
	Usage    string // help message
	Value    Value  // value as set
	DefValue string // default value (as text); for usage message
}
~~~

与之相关的是 `FlagSet` 结构体：

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

有了这部分知识，我们再从 flag.go 文件中定义的全局变量和 `init()` 说起：

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

func commandLineUsage() {
	Usage()
}

// NewFlagSet returns a new, empty flag set with the specified name and
// error handling property. If the name is not empty, it will be printed
// in the default usage message and in error messages.
func NewFlagSet(name string, errorHandling ErrorHandling) *FlagSet {
	f := &FlagSet{
		name:          name,
		errorHandling: errorHandling,
	}
	f.Usage = f.defaultUsage
	return f
}
~~~

`CommandLine` 是默认的命令行参数 `FlagSet`，用于解析 `os.Args`。创建的 `FlagSet` 变量，其 name 设置为 `os.Args[0]`。

现在，我们以 `int` 作为例子，说明整个注册、解析过程！再次之前先看看 `int` 为此打下前提：

~~~go
// -- int Value
type intValue int

func newIntValue(val int, p *int) *intValue {
	*p = val
	return (*intValue)(p) // 强制类型转换
}

func (i *intValue) Set(s string) error
func (i *intValue) Get() interface{}
func (i *intValue) String() string
~~~

我们之前提到的 flag 包支持 8 种类型的命令行参数解析，flag.go 中均有相类似的方法定义。

命令行参数==注册的过程==，调用了：

~~~go
// Int defines an int flag with specified name, default value, and usage string.
// The return value is the address of an int variable that stores the value of the flag.
func Int(name string, value int, usage string) *int {
	return CommandLine.Int(name, value, usage)
}

// Int defines an int flag with specified name, default value, and usage string.
// The return value is the address of an int variable that stores the value of the flag.
func (f *FlagSet) Int(name string, value int, usage string) *int {
    // CommandLine 作为 FlagSet
	p := new(int)
	f.IntVar(p, name, value, usage)
	return p
}

// IntVar defines an int flag with specified name, default value, and usage string.
// The argument p points to an int variable in which to store the value of the flag.
func (f *FlagSet) IntVar(p *int, name string, value int, usage string) {
	f.Var(newIntValue(value, p), name, usage)
}

// Var defines a flag with the specified name and usage string. The type and
// value of the flag are represented by the first argument, of type Value, which
// typically holds a user-defined implementation of Value. For instance, the
// caller could create a flag that turns a comma-separated string into a slice
// of strings by giving the slice the methods of Value; in particular, Set would
// decompose the comma-separated string into the slice.
func (f *FlagSet) Var(value Value, name string, usage string) {
	// Remember the default value as a string; it won't change.
	flag := &Flag{name, usage, value, value.String()} // 创建 Flag 变量，并初始化
	_, alreadythere := f.formal[name]
	if alreadythere {
		var msg string
		if f.name == "" {
			msg = fmt.Sprintf("flag redefined: %s", name)
		} else {
			msg = fmt.Sprintf("%s flag redefined: %s", f.name, name)
		}
		fmt.Fprintln(f.Output(), msg)
		panic(msg) // Happens only if flags are declared with identical names
	}
	if f.formal == nil {
		f.formal = make(map[string]*Flag)
	}
    // 把注册进来的 Flag 添加到 f.formal 中，此处 f 是指 CommandLine
	f.formal[name] = flag
}

// IntVar defines an int flag with specified name, default value, and usage string.
// The argument p points to an int variable in which to store the value of the flag.
func IntVar(p *int, name string, value int, usage string) {
	CommandLine.Var(newIntValue(value, p), name, usage)
}

// Var defines a flag with the specified name and usage string. The type and
// value of the flag are represented by the first argument, of type Value, which
// typically holds a user-defined implementation of Value. For instance, the
// caller could create a flag that turns a comma-separated string into a slice
// of strings by giving the slice the methods of Value; in particular, Set would
// decompose the comma-separated string into the slice.
func Var(value Value, name string, usage string) {
	CommandLine.Var(value, name, usage)
}
~~~

注册的方式可以归纳为：

1. `Type()`：比如 `Int()` 等；
2. `TypeVar()`：比如 `IntVar()` 等。相对于 `Type()` 更简洁的原因是，其 `TypeVar()` 的参数 `p *int`；
3. `Var`：专用于自定义类型的注册。

注册过程执行完后，`f.formal[name] = flag`，及将 Flag 对象添加到了 `FlagSet` 的 Map 类型成员中；而 `Flag` 成员中，最重要的是 Value。

命令行==解析的过程==，调用了：

~~~go
// Parse parses the command-line flags from os.Args[1:]. Must be called
// after all flags are defined and before flags are accessed by the program.
func Parse() {
	// Ignore errors; CommandLine is set for ExitOnError.
	CommandLine.Parse(os.Args[1:])
}

// Parse parses flag definitions from the argument list, which should not
// include the command name. Must be called after all flags in the FlagSet
// are defined and before flags are accessed by the program.
// The return value will be ErrHelp if -help or -h were set but not defined.
func (f *FlagSet) Parse(arguments []string) error {
	f.parsed = true
	f.args = arguments
	for {
		seen, err := f.parseOne() // 核心调用
		if seen {
			continue
		}
		if err == nil {
			break
		}
		switch f.errorHandling {
		case ContinueOnError:
			return err
		case ExitOnError:
			if err == ErrHelp {
				os.Exit(0)
			}
			os.Exit(2)
		case PanicOnError:
			panic(err)
		}
	}
	return nil
}

// parseOne parses one flag. It reports whether a flag was seen.
func (f *FlagSet) parseOne() (bool, error) {
	if len(f.args) == 0 {
		return false, nil
	}
	s := f.args[0]
	if len(s) < 2 || s[0] != '-' {
		return false, nil
	}
	numMinuses := 1
	if s[1] == '-' {
		numMinuses++
		if len(s) == 2 { // "--" terminates the flags
			f.args = f.args[1:]
			return false, nil
		}
	}
	name := s[numMinuses:]
	if len(name) == 0 || name[0] == '-' || name[0] == '=' {
		return false, f.failf("bad flag syntax: %s", s)
	}

	// it's a flag. does it have an argument?
	f.args = f.args[1:] // 循环遍历，舍弃已经解析的命令行参数
	hasValue := false
	value := ""
	for i := 1; i < len(name); i++ { // equals cannot be first
		if name[i] == '=' {
            // 取到 name-value
			value = name[i+1:]
			hasValue = true
			name = name[0:i]
			break
		}
	}
    // 取出原先注册的 map[string]*Flag
	m := f.formal
	flag, alreadythere := m[name] // BUG
	if !alreadythere {
		if name == "help" || name == "h" { // special case for nice help message.
			f.usage()
			return false, ErrHelp
		}
		return false, f.failf("flag provided but not defined: -%s", name)
	}

    // 对取到的 Flag 遍历做类型判断
	if fv, ok := flag.Value.(boolFlag); ok && fv.IsBoolFlag() { // special case: doesn't need an arg
		if hasValue {
			if err := fv.Set(value); err != nil {
				return false, f.failf("invalid boolean value %q for -%s: %v", value, name, err)
			}
		} else {
			if err := fv.Set("true"); err != nil {
				return false, f.failf("invalid boolean flag %s: %v", name, err)
			}
		}
	} else {
		// It must have a value, which might be the next argument.
		if !hasValue && len(f.args) > 0 {
			// value is the next arg
			hasValue = true
			value, f.args = f.args[0], f.args[1:]
		}
		if !hasValue {
			return false, f.failf("flag needs an argument: -%s", name)
		}
        // 实际调用的是 *intValue 作为接收者的 Set 方法
		if err := flag.Value.Set(value); err != nil {
			return false, f.failf("invalid value %q for flag -%s: %v", value, name, err)
		}
	}
	if f.actual == nil {
		f.actual = make(map[string]*Flag)
	}
    // 解析完成后，设置 FlagSet 结构体中的 actual 成员
	f.actual[name] = flag
	return true, nil
}

func (i *intValue) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, strconv.IntSize)
	if err != nil {
		err = numError(err)
	}
    // 更新 Value 值
	*i = intValue(v)
	return err
}
~~~









# 框架设计思想



