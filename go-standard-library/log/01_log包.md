Go 中的 log 标准库实现了简单的日志功能封装，其中定义了 Logger 类型，并附带有格式化输出的方法。在 log 包中，有一个预定义的，名为 std 的 Logger 实例，可直接使用 std 这个实例。

~~~go
var std = New(os.Stderr, "", LstdFlags)

// Default returns the standard logger used by the package-level output functions.
func Default() *Logger { return std }
~~~

如果直接调用 log 标准库下的方法，比如：

~~~go
// Print calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Print.
func Print(v ...interface{}) {
	std.Output(2, fmt.Sprint(v...))
}

// Fatal is equivalent to Print() followed by a call to os.Exit(1).
func Fatal(v ...interface{}) {
	std.Output(2, fmt.Sprint(v...))
	os.Exit(1)
}

// Panic is equivalent to Print() followed by a call to panic().
func Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	std.Output(2, s)
	panic(s)
}
~~~

则使用的就是 std 这个 Logger 实例。从实例的定义中可以看出，会将日志输出到 os.Stderr 中，同时，在每条日志前都会附带上日期和时间。

# 1 Logger 结构

~~~go
// A Logger represents an active logging object that generates lines of
// output to an io.Writer. Each logging operation makes a single call to
// the Writer's Write method. A Logger can be used simultaneously from
// multiple goroutines; it guarantees to serialize access to the Writer.
type Logger struct {
	mu     sync.Mutex // ensures atomic writes; protects the following fields
	prefix string     // prefix on each line to identify the logger (but see Lmsgprefix)
	flag   int        // properties
	out    io.Writer  // destination for output
	buf    []byte     // for accumulating text to write
}
~~~

Logger 代表的是一个活动的一个 Logger 实例，可将内容输出到指定的 io.Writer 中。每调用依次 Logger 的方法，都会调用 io.Writer 的 write 方法。Logger 可以在多个 goroutine 中并发使用，会序列化的访问 io.Writer 实例。

# 2 创建实例

可以使用 New 函数创建一个 Logger 实例：

~~~go
// New creates a new Logger. The out variable sets the
// destination to which log data will be written.
// The prefix appears at the beginning of each generated log line, or
// after the log header if the Lmsgprefix flag is provided.
// The flag argument defines the logging properties.
func New(out io.Writer, prefix string, flag int) *Logger {
	return &Logger{out: out, prefix: prefix, flag: flag}
}
~~~

来看看 log 标准库中是如何使用该函数创建默认实例的：

~~~go
var std = New(os.Stderr, "", LstdFlags)

// Default returns the standard logger used by the package-level output functions.
func Default() *Logger { return std }
~~~



