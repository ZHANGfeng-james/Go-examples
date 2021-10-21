Go 中的 flag 包，其用途是**解析 command-line 标识符的值**。

> 解析：parse

# 1 使用方法

意思就是说，**解析**如下的**命令行标识符**或者称为**命令行参数**：

~~~go
go run main.go -name 1
~~~

在运行程序中获取到 Command Line 输入的参数 `-name`，经过解析得到了 key-value，即为：`name=1`。

> 如果命令行中输入参数，但是程序中**并没有定义该参数**。此时应用程序会提示：`flag provided but not defined: -name` 意思就是说，**程序中并没有定义**。**程序不会正常运行！**

使用 flag 包**定义标识符**，可以使用 flag.String()、flag.Bool() 或者 flag.Int() 等。

如下声明了一个名为 n 的 int 类型标识符，并将解析后的值存在 `nFlag` 变量中，其类型是 `*int`：

~~~go
func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var nFlag = flag.Int("n", 0, "flag param n value") // 在程序中定义标识符，名称为 n

	if !flag.Parsed() {
		flag.Parse()
	}
	log.Printf("nFlag:%d", *nFlag) // 使用 log 包功能输出日志
}
~~~

`flag.Int` 函数包含 3 个参数，其中默认值 defVal 表示如果在命令行中并没有指定该参数，则该值就是默认的值。

同样，可以使用如下这种，**将命令行参数绑定到指定的变量上**：

~~~go
func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var nFlag = flag.Int("n", 0, "flag param n value")

	var nameFlag string
	flag.StringVar(&nameFlag, "name", "Katyusha", "name of company")

	if !flag.Parsed() {
		flag.Parse()
	}
	log.Printf("nFlag:%d, nameFlag:%s", *nFlag, nameFlag)
}
~~~

在上述定义的 int 和 string 类型的命令行参数中，可有如下**命令行输入形式**，均能让 flag 解析出给定的参数值：

1. `go run usage.go -n 1 -name Arthur`
2. `go run usage.go -n 1 -name "Arthur"`
3. `go run usage.go -n 1 -name="Arthur"` 同样的，对 int 类型的参数也是一样的。

**对于 bool 类型的参数**，有不同输入形式：

~~~go
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var nFlag = flag.Int("i", 0, "flag param n value")

	var nameFlag string
	flag.StringVar(&nameFlag, "s", "Katyusha", "name of company")

	var boolFlag bool
	flag.BoolVar(&boolFlag, "b", false, "whether or not") // bool 类型的

	if !flag.Parsed() {
		flag.Parse()
	}
	log.Printf("nFlag:%d, nameFlag:%s, boolFlag:%v", *nFlag, nameFlag, boolFlag)
}
~~~

输入命令：`go run usage.go -b` 就表示已给定了参数，其参数值为 true，不需要在其后加上 true 值！

或者，也可以**创建自定义的标识结构体类型**，同时实现 Value 接口，即可以创建自定义结构的命令行标识符：

~~~go
type CarInfo struct {
	brand string
}

func (carInfo *CarInfo) String() string {
	return carInfo.brand
}

func (carInfo *CarInfo) Set(value string) error {
	carInfo.brand = value
	return nil
}

var car CarInfo
flag.Var(&car, "carinfo", "param car info")
~~~

如果在进程中**定义了某个标识符**，但是实际运行时并没有在 Command Line 中输入该标识符，此时进程中与该标识符绑定的变量值**就是其默认值**。

在所有标识符都定义完后，接下来就可以调用 `flag.Parse()`，即为：**解析**命令行参数，并**赋值**到对应的变量中。

在调用 `flag.Parse()` 后，如果还需要附带有**其他参数**，可以**在标识符参数之后添加**，其获取方式：

~~~go
func main() {
	flag.Parse()
	fmt.Println(flag.Args())
}
PS G:\michoi> go run main.go 123 true
[123 true]
~~~

**P.S.**

一定记得要在命名行标识符定义完成**之后**再调用 `flag.Parse()`，以此触发命令行参数的解析。其作用相当于是，让 flag 包的程序预先知道应用程序需要解析哪些命令行参数（这就是**定义命令行标识符**的含义）。

# 2 命令行参数语法

允许如下语法：

~~~bash
-flag
-flag=x
-flag x  // non-boolean flags only，bool 类型的参数值不允许这种形式
~~~

也就是说，非 bool 类型可以使用上述 3 种形式，bool 类型仅能使用前 2 种形式。**可以使用 `-` 或者 `--` 形式**，是等价的。

对于整型值的标识符参数，允许 1234、0664、0x1234 以及负数值，对应的就是十进制、八进制和十六进制格式；**对于 bool 类型值**，允许：

~~~bash
1, 0, t, f, T, F, true, false, TRUE, FALSE, True, False
~~~

也就是说，对于 bool 类型的值，使用上述所有的值都是可以为 bool 类型的命令行参数赋值的。

对于 Duration 类型的标识符参数，允许的值是从 `time.ParseDuration` 返回的值。

# 3 Flag 示例程序

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
	flag.StringVar(&gopherType, "g", defaultGopher, usage+" (sharthand)")
}

var species = flag.String("species", "gopher", "the species we are studying")

func main() {
	flag.Parse()

	fmt.Println(gopherType, ", ", *species)
}
PS G:\GoUsage> go run main.go -gopher_type="Rocket" -g="rocket"
rocket ,  gopher
PS G:\GoUsage> go run main.go -gopher_type="Rocket" -g="rocket" -species "Gopher"
rocket ,  Gopher
~~~

**自定义类型**的标识符参数解析：

~~~go
package main

import (
	"errors"
	"flag"
	"fmt"
	"strings"
	"time"
)

type interval []time.Duration

func (i *interval) String() string {
	return fmt.Sprint(*i)
}

func (i *interval) Set(value string) error {
	if len(*i) > 0 {
		return errors.New("interval flag alread set")
	}

	for _, dt := range strings.Split(value, ",") { // 按照输入的格式解析参数
		duration, err := time.ParseDuration(dt)
		if err != nil {
			return err
		}
		*i = append(*i, duration)
	}
	return nil
}

var intervalFlag interval

func init() {
	flag.Var(&intervalFlag, "deltaT", "comma-separated list of intervals to use between events")
}

func main() {
	flag.Parse()
	fmt.Println(intervalFlag.String())
}
PS G:\GoUsage> go run main.go -deltaT 10s,15s
[10s 15s]
~~~

`func (i *interval) Set(value string)` 的入参，如果给定的命令行是 `go run main.go -deltaT 10s,15s`，则 value 值就是：`10s,15s`。

# 4 FlagSet 的使用示例

我们来看一个最简单的实例：

~~~go
package main

import (
	"flag"
	"fmt"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// 实际的作用就是取到 flag.Args，使用默认的 FlagSet
	flag.Parse()

	var name string

	goCmd := flag.NewFlagSet("go", flag.ExitOnError)       // 创建 name 为 go 的 FlagSet
	goCmd.StringVar(&name, "name", "Go语言", "help message") // goCmd 这个 FlagSet 中为 name 变量预解析参数标识符

	javaCmd := flag.NewFlagSet("java", flag.ExitOnError)
	javaCmd.StringVar(&name, "name", "Java语言", "help message") // javaCmd 这个 FlagSet 中为 name 变量预解析参数标识符

	// 取到 os.Args[1:] 的命令行参数
	args := flag.Args()
	if len(args) <= 0 {
		return
	}
	fmt.Printf("%d, %v\n", len(args), args)

	// 匹配到对应的子命令
	switch args[0] {
	case "go":
		// 解析接下来的命令行参数
		_ = goCmd.Parse(args[1:])
	case "java":
		_ = javaCmd.Parse(args[1:])
	}

	fmt.Printf("name=%q\n", name)
}
~~~

这就是 FlagSet 的另一个用途：**为指定的名称创建一个 Flag 集合**。实际上 `flag.Parse()` 在上述示例程序中只有一个作用，就是获取到 `flag.Args[]`。在使用 FlagSet 中，最重要的是 `flagSet.Parse(args[1:])`，**触发这个 flagSet 解析命令行参数**。对应的命令行输入内容是：

~~~bash
ant@MacBook-Pro v1 % go run usage.go go -name=Katyusha
2, [go -name=Katyusha]
name="Katyusha"
~~~

其中 go 是一个 FlagSet，其解析的开始部分就是 `args[1:]`，也就是根据命令行参数的 `go` 内容，**作为一个起始部分开始解析这个指定的 FlagSet 对应的 Flag**。

> go 或者是 java 是一个 FlagSet，其等价于在 Flag 示例程序中以进程名命名的 FlaSet。在使用时，不同之处是：触发 FlagSet 执行解析参数的入参是不同的。对于 go 这样的 FlagSet 需要自行给定入参，比如：`goCmd.Parse(args[1:])`

# 5 其他方式获取参数

通过前面，我们已经知道了 Command Line 标识符参数的获取、解析方式，这是一种获取参数的方法。

**非标识符参数**的获取方式和入参形式：

~~~go
package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
    flag.Parse() // 必须要调用 flag.Parse()

	flags := flag.Args()
	fmt.Println(flags)
}
PS G:\GoUsage> go run main.go 123 true
[123 true]
~~~

对于**非标识符参数**，也就是**在标识符参数之后**附加的一些值，可直接通过 `flag.Args()` 获取到。`flag.Args()` 的返回结果是在命令行参数解析完成之后剩下的部分：

~~~go
ant@MacBook-Pro v1 % go run usage.go -i 20 -s michoi -b true 1
2021/10/16 11:02:08 usage.go:24: nFlag:20, nameFlag:michoi, boolFlag:true
2021/10/16 11:02:08 usage.go:27: [/var/folders/3w/0pprfk1x13lfrt7vdxryd75r0000gn/T/go-build460835673/b001/exe/usage -i 20 -s michoi -b true 1]
2021/10/16 11:02:08 usage.go:30: [true 1]
~~~

`[true 1]` 就是 `flag.Parse()` 指定的**剩下部分**。

还有一种方式也能够获取到**非标识符参数**：

~~~go
package main

import (
	"fmt"
	"os"
)

func main() {
	all := os.Args
	fmt.Printf("size:%d, Parameters: %v\n", len(os.Args), all[1:])
}
PS G:\GoUsage> go run main.go 123 true
size:3, Parameters: [123 true]
~~~

`os.Args` 获取到的是所有参数内容，比如：命令行参数和非选项参数。其返回值类型是 `[]string`，首个元素是应用程序进程，紧接着的就是所有的参数。比如：

~~~go
ant@MacBook-Pro v1 % go run usage.go -i 20 -s michoi -b true 1
2021/10/16 10:58:54 usage.go:24: nFlag:20, nameFlag:michoi, boolFlag:true
2021/10/16 10:58:54 usage.go:27: [/var/folders/3w/0pprfk1x13lfrt7vdxryd75r0000gn/T/go-build805321147/b001/exe/usage -i 20 -s michoi -b true 1]
~~~

本质上来说，标准库 flag 的命令行参数解析的数据来源也是 `os.Args`。

# 6 命令行参数解析原理

回答如下疑问：

1. 为什么会出现类似 `flag provided but not defined: -name` 错误（**程序不会正常运行**），如何产生第 6 章的 2 种错误的？
2. flag 标准库到底是如何解析出命令行参数的？

从最简单的 flag 标准库的使用 Demo 开始分析：

~~~go
package main

import (
	"flag"
	"fmt"
)

func init() {
	name := flag.String("name", "Katyusha", "name of action")

	var age int
	flag.IntVar(&age, "age", 18, "age of person")

	var inSchool bool
	flag.BoolVar(&inSchool, "inschool", false, "whether in school")

	flag.Parse()

	fmt.Printf("%q: %d - %v\n", *name, age, inSchool)
}

func main() {

}
~~~

上述 Demo，首先引入了 flag 标准库，在该标准库中做了一些初始化工作：

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

首先需要**弄清楚的是 2 个概念**：

1. **FlagSet**：可以看做是**一个标识符的集合**，这个集合对应的结构体就是 FlagSet 的各个字段内容，比如 args 就是对应的参数。
2. **CommandLine 变量**：默认的 FlagSet。从其初始化来看，使用了 `os.Args[0]` 这个名字作为 FlagSet 的名称。而且 `os.Args[0]` 对应的就是当前运行进程的绝对路径，比如 `C:\Users\ADMINI~1\AppData\Local\Temp\go-build3140955443\b001\exe\main.exe`。

依据初始化内容，我们实际上可以看出 CommandLine 默认构造的就是以自己进程名的 FlagSet，这个实例时默认创建的。依据 Demo 程序，紧接着的就是**预解析（绑定）**对应名称的**标识符参数**：

~~~go
// String defines a string flag with specified name, default value, and usage string.
// The return value is the address of a string variable that stores the value of the flag.
func String(name string, value string, usage string) *string {
	return CommandLine.String(name, value, usage)
}

// String defines a string flag with specified name, default value, and usage string.
// The return value is the address of a string variable that stores the value of the flag.
func (f *FlagSet) String(name string, value string, usage string) *string {
	p := new(string)
	f.StringVar(p, name, value, usage)
	return p
}

// StringVar defines a string flag with specified name, default value, and usage string.
// The argument p points to a string variable in which to store the value of the flag.
func (f *FlagSet) StringVar(p *string, name string, value string, usage string) {
    // 将定义的 Flag 加入到默认的 FlaSet 中去，相当于是内存缓存
	f.Var(newStringValue(value, p), name, usage)
}

func newStringValue(val string, p *string) *stringValue {
    // 设置默认值
	*p = val
	return (*stringValue)(p)
}

// -- string Value
type stringValue string

// Value is the interface to the dynamic value stored in a flag.
// (The default value is represented as a string.)
//
// If a Value has an IsBoolFlag() bool method returning true,
// the command-line parser makes -name equivalent to -name=true
// rather than using the next command-line argument.
//
// Set is called once, in command line order, for each flag present.
// The flag package may call the String method with a zero-valued receiver,
// such as a nil pointer.
type Value interface {
	String() string
	Set(string) error
}

func (s *stringValue) Set(val string) error {
    // stringValue 类型
	*s = stringValue(val)
	return nil
}

func (s *stringValue) Get() interface{} { return string(*s) }

func (s *stringValue) String() string { return string(*s) }

// Var defines a flag with the specified name and usage string. The type and
// value of the flag are represented by the first argument, of type Value, which
// typically holds a user-defined implementation of Value. For instance, the
// caller could create a flag that turns a comma-separated string into a slice
// of strings by giving the slice the methods of Value; in particular, Set would
// decompose the comma-separated string into the slice.
func (f *FlagSet) Var(value Value, name string, usage string) {
	// Remember the default value as a string; it won't change.
	flag := &Flag{name, usage, value, value.String()}
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
    // 每次预解析一个 flag，都会将 name 加入到 f.formal 中，作为内存缓存
	f.formal[name] = flag
}

// A Flag represents the state of a flag.
type Flag struct {
	Name     string // name as it appears on command line
	Usage    string // help message
	Value    Value  // value as set
	DefValue string // default value (as text); for usage message
}
~~~

很明显，使用的是 CommandLine 这个 FlagSet 实例，去解析指定名称的命令行标识符参数值。

如果是 int 类型，不同的部分是这些：

~~~go
// -- int Value
type intValue int

func newIntValue(val int, p *int) *intValue {
    // 设置默认值
	*p = val
	return (*intValue)(p)
}

func (i *intValue) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, strconv.IntSize)
	if err != nil {
		err = numError(err)
	}
    // intValue 类型
	*i = intValue(v)
	return err
}

func (i *intValue) Get() interface{} { return int(*i) }

func (i *intValue) String() string { return strconv.Itoa(int(*i)) }
~~~

如果是 bool 类型，不同的部分是这些：

~~~go
// -- bool Value
type boolValue bool

func newBoolValue(val bool, p *bool) *boolValue {
    // 设置默认值
	*p = val
	return (*boolValue)(p)
}

func (b *boolValue) Set(s string) error {
	v, err := strconv.ParseBool(s)
	if err != nil {
		err = errParse
	}
    // boolValue 类型
	*b = boolValue(v)
	return err
}

func (b *boolValue) Get() interface{} { return bool(*b) }

func (b *boolValue) String() string { return strconv.FormatBool(bool(*b)) }
~~~

接下来就是，就是真正的解析部分了：

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
		seen, err := f.parseOne()
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
    // 分析了 name 前的 - 或者 -- 标识，拆分出 name；此处明显看到仅支持 - 或者 -- 作为前缀
	name := s[numMinuses:]
	if len(name) == 0 || name[0] == '-' || name[0] == '=' {
		return false, f.failf("bad flag syntax: %s", s)
	}

	// it's a flag. does it have an argument?
	f.args = f.args[1:]
	hasValue := false
	value := ""
	for i := 1; i < len(name); i++ { // equals cannot be first
		if name[i] == '=' {
			value = name[i+1:]
			hasValue = true
			name = name[0:i]
			break
		}
	}
	m := f.formal
    // 从 FlagSet 中，依据 name 取到对应的 Flag 实例
	flag, alreadythere := m[name] // BUG
	if !alreadythere {
		if name == "help" || name == "h" { // special case for nice help message.
			f.usage()
			return false, ErrHelp
		}
        // 我们看到了熟悉的身影！
		return false, f.failf("flag provided but not defined: -%s", name)
	}

    // 依据 flag，取到对应的 Value 值，判断当前是否是 boolFlag
	if fv, ok := flag.Value.(boolFlag); ok && fv.IsBoolFlag() { // special case: doesn't need an arg
		if hasValue {
            // 如果是类似 -isLogin=true，调用 *boolValue 类型的 Set 方法，其入参就是 value
			if err := fv.Set(value); err != nil {
				return false, f.failf("invalid boolean value %q for -%s: %v", value, name, err)
			}
		} else {
            // 如果是类似 -isLogin，默认就标识设置了该标志位
			if err := fv.Set("true"); err != nil {
				return false, f.failf("invalid boolean flag %s: %v", name, err)
			}
		}
	} else {
		// It must have a value, which might be the next argument.
		if !hasValue && len(f.args) > 0 {
			// value is the next arg
			hasValue = true
            // 取到下一个参数
			value, f.args = f.args[0], f.args[1:]
		}
		if !hasValue {
			return false, f.failf("flag needs an argument: -%s", name)
		}
        // 调用对应 flag.Value 类型的 Set 方法，入参是 value
		if err := flag.Value.Set(value); err != nil {
			return false, f.failf("invalid value %q for flag -%s: %v", value, name, err)
		}
	}
	if f.actual == nil {
		f.actual = make(map[string]*Flag)
	}
    // 将解析完成后的 name，填充到的 actual 中
	f.actual[name] = flag
	return true, nil
}
~~~

我们从上述整个过程看到了标识符参数的**预解析（绑定）** formal，还看到真实的**解析**过程 actual。

从概念范围来看，从大到小依次是：FlagSet -- Flag -- Value，也就是标识符集合（比如 main.go）、标识符（比如 -name）、具体的参数值（比如 -name="Katyusha"）。

上面是从命令行作为输入端，也就是说在 CLI 中输入对应的 `-flagname`，给定值，经过 Parse 后得到的就是指定的输入值。那**反其道而行之**：在绑定变量后，直接赋值，会有什么效果？

~~~go
func ReadFromVariable() {
	var nFlag = flag.Int("i", 10, "flag param n value")
	*nFlag = 20

	flag := flag.CommandLine.Lookup("i")
	log.Printf("read:%s", flag.Value.String())
}
~~~

也就是说，如果直接给 Flag 绑定的变量赋值，（不需要 Parse 的情况下）直接就可以从 Flag 中获取到值。**原因：`*nFlag` 指向的变量已经和 flag 绑定在一起了**！

# 7 错误解决

## 7.1 flag provided but not defined

在初次使用时，遇到了错误：

~~~bash
PS G:\GoUsage> go run main.go -name 1
line 1
flag provided but not defined: -name
Usage of C:\Users\DEVELO~1\AppData\Local\Temp\go-build956667662\b001\exe\main.exe:
~~~

错误信息很明显：`flag provided but not defined: -name`

源代码如下：

~~~go
package main

import (
	"flag"
	"fmt"
)

func main() {
	var nFlag int
	flag.Parse()
	fmt.Println("name:", nFlag)
}
~~~

源代码是很简单的，但就是在 `flag.Parse()` 中报错！原因是：**程序在执行 `flag.Parse()` 检测到了 `-name` 这个标志位，但是进程中并没有定义这个名为 `name` 的标识符（这种情况不被允许）** 。反之，如果在执行时，没有给定已在程序解析的标识符，此时标识符对应的变量就被赋予其默认值，也就是说，这种情况是**被允许的**。因此，如果想要在进程中解析出**命令行中已给出的标识符**，就应该在进程中定义该标识符，比如定义名为 name 这个标识符。否则，程序会报错。

## 7.2 flag 没有正常被解析

还有另外一种情况：

~~~go
package main

import (
	"flag"
	"fmt"
)

func init() {
	flag.Parse()

	name := flag.String("name", "Katyusha", "name of action")

	var age int
	flag.IntVar(&age, "age", 18, "age of person")

	var inSchool bool
	flag.BoolVar(&inSchool, "inschool", false, "whether in school")

	fmt.Printf("%q: %d - %v\n", *name, age, inSchool)
}

func main() {

}
~~~

程序执行过程中出现如下错误：

~~~bash
PS E:\go_developer_roadmap\ProgrammingLanguage\Go Standard Interface\GoUsage> go run main.go -name="Qru" -age=3 -inschool
flag provided but not defined: -name
Usage of C:\Users\ADMINI~1\AppData\Local\Temp\go-build171152850\b001\exe\main.exe:
exit status 2
~~~

也就是说，flag 标准库无法检查到程序中已注册的命令行标识符 -name 标识，和第一种情况是相同的报错提示。

错误的原因是，`flag.Parse()` 函数调用顺序错误！**正确的调用顺序**是：

~~~go
package main

import (
	"flag"
	"fmt"
)

func init() {
	name := flag.String("name", "Katyusha", "name of action")

	var age int
	flag.IntVar(&age, "age", 18, "age of person")

	var inSchool bool
	flag.BoolVar(&inSchool, "inschool", false, "whether in school")

    flag.Parse()
	fmt.Printf("%q: %d - %v\n", *name, age, inSchool)
}

func main() {

}
~~~



