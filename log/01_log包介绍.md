Go 中的 log 标准库实现了简单的日志功能封装，其中定义了 Logger 类型，并附带有格式化输出的方法。在 log 包中，有一个预定义的，**名为 std 的 Logger 实例**，可直接使用 std 这个实例：

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

需要注意的点是：

1. 可以使用的导出方法是：Print、Fatal 和 Panic，以及与之相关的方法。这些导出方法是在 std 这个 Logger 实例上执行的，日志信息被输出到 `os.Stderr`——**标准错误输出**中；
2. Fatal 在输出日志后会**调用 os.Exit(1)**；
3. Panic 在输出日志后会**调用 panic 从而引发 Panic**。

总体来看，**log 标准库 API 分为 3 类**：

* 基于 std —— 引入 log 包时创建的 Logger 实例，有固定的格式，没有 prefix——的可导出函数；
* 定义的 Logger 结构体，并在其上定义的方法（*Logger 类型）；
* 可导出的方法，让用户创建 *Logger，比如 New 函数。

### 1 Logger 结构

~~~go
// A Logger represents an active logging object that generates lines of
// output to an io.Writer. Each logging operation makes a single call to
// the Writer's Write method. A Logger can be used simultaneously from
// multiple goroutines; it guarantees to serialize access to the Writer.
type Logger struct {
	mu     sync.Mutex // ensures atomic writes; protects the following fields
	prefix string     // prefix on each line to identify the logger (but see Lmsgprefix)
    
	flag   int        // properties 属性值，区分Logger的功能
    
	out    io.Writer  // destination for output
	buf    []byte     // for accumulating text to write 输出缓冲区，迟早会被消费的内容
}
~~~

Logger 代表的是一个活动的一个 Logger 实例，可将内容输出到指定的 io.Writer 中。每调用依次 Logger 的方法，都会调用 io.Writer 的 write 方法。Logger 可以在**多个 goroutine 中并发使用**，会序列化的访问 io.Writer 实例（确保**并发安全**）。

### 2 创建实例

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

### 3 Logger 属性值

~~~go
func usePrefix() {
	var helloStr = "Hello, log!"
	stdFlags := log.Flags() // LstdFlags = Ldate | Ltime
	log.Print(helloStr, " stdFlags:", stdFlags)

	log.SetPrefix("---")

	var prefixFlag map[int]string = make(map[int]string)
	prefixFlag[log.LUTC] = "log.LUTC"
	prefixFlag[log.Ldate] = "log.Ldate"
	prefixFlag[log.Ltime] = "log.Ltime"
	prefixFlag[log.Lmicroseconds] = "log.Lmicroseconds"
	prefixFlag[log.LstdFlags] = "log.LstdFlags"
	prefixFlag[log.Llongfile] = "log.Llongfile"
	prefixFlag[log.Lshortfile] = "log.Lshortfile"
	prefixFlag[log.Lmsgprefix] = "log.Lmsgprefix"

	for prefix, str := range prefixFlag {
		log.SetFlags(prefix)
		fmt.Printf("%+20s:", str)
		log.Print(helloStr)
	}
}

ant@MacBook-Pro v1 % go run log.go
2021/10/15 20:46:01 Hello, log! stdFlags:3
      log.Lmsgprefix:---Hello, log!
            log.LUTC:---Hello, log!
           log.Ldate:---2021/10/15 Hello, log!
           log.Ltime:---20:46:01 Hello, log!
   log.Lmicroseconds:---20:46:01.133758 Hello, log!
       log.LstdFlags:---2021/10/15 20:46:01 Hello, log!
       log.Llongfile:---/Users/ant/Documents/ProgrammingLanguage/Go/Go-examples-with-tests/log/v1/log.go:30: Hello, log!
      log.Lshortfile:---log.go:30: Hello, log!
~~~

属性值 Properties（在结构体中**用 flag 存储**） 用于标记 Logger **输出的内容及其形式**。

**P.S.**

1. log.Lmsgprefix：被置位时，原先放在 line 前面的内容，会被放置到 message 前面；
2. 如果 log.Lmsgprefix 没有被置位，则会放置在 line 前面。

### 4 打印输出日志

打印输出调用类似 Print 的方法：

~~~go
// Printf calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Printf(format string, v ...interface{}) {
	l.Output(2, fmt.Sprintf(format, v...))
}
~~~

作为**一次 Logger Event**，最终汇集到：

~~~go
// Output writes the output for a logging event. The string s contains
// the text to print after the prefix specified by the flags of the
// Logger. A newline is appended if the last character of s is not
// already a newline. Calldepth is used to recover the PC and is
// provided for generality, although at the moment on all pre-defined
// paths it will be 2.
func (l *Logger) Output(calldepth int, s string) error {
	now := time.Now() // get this early.
	var file string
	var line int
    
	l.mu.Lock() // 加锁，确保并发安全
	defer l.mu.Unlock()
    
	if l.flag&(Lshortfile|Llongfile) != 0 { // Lshortfile|Llongfile才需要源文件名
		// Release lock while getting caller info - it's expensive.
        l.mu.Unlock() // Unlock() 之后是否会导致下面内容出现竞争？
        
		var ok bool
		_, file, line, ok = runtime.Caller(calldepth) // 确定的某一层 calldepth 的栈帧
		if !ok {
			file = "???"
			line = 0
		}
		l.mu.Lock()
	}
    
	l.buf = l.buf[:0] // clean up
    
	l.formatHeader(&l.buf, now, file, line)
	l.buf = append(l.buf, s...)
	if len(s) == 0 || s[len(s)-1] != '\n' {
		l.buf = append(l.buf, '\n')
	}
	_, err := l.out.Write(l.buf) // 输出到 l.out
	return err // ？？？返回值有什么用？
}
~~~

基于整个的流程，在获取到**前置信息**后，需要转化成 formatHeader：**严格的输出顺序**

~~~go
func (l *Logger) formatHeader(buf *[]byte, t time.Time, file string, line int) {
	if l.flag&Lmsgprefix == 0 {
        // l.prefix (if it's not blank and Lmsgprefix is unset)
		*buf = append(*buf, l.prefix...)
	}
    // date and/or time (if corresponding flags are provided)
	if l.flag&(Ldate|Ltime|Lmicroseconds) != 0 {
		if l.flag&LUTC != 0 {
            t = t.UTC() // t 被UTC()覆写
		}
		if l.flag&Ldate != 0 {
			year, month, day := t.Date()
			itoa(buf, year, 4) // integer to ascii
			*buf = append(*buf, '/')
			itoa(buf, int(month), 2)
			*buf = append(*buf, '/')
			itoa(buf, day, 2)
			*buf = append(*buf, ' ')
		}
		if l.flag&(Ltime|Lmicroseconds) != 0 {
			hour, min, sec := t.Clock()
			itoa(buf, hour, 2)
			*buf = append(*buf, ':')
			itoa(buf, min, 2)
			*buf = append(*buf, ':')
			itoa(buf, sec, 2)
			if l.flag&Lmicroseconds != 0 {
				*buf = append(*buf, '.')
				itoa(buf, t.Nanosecond()/1e3, 6)
			}
			*buf = append(*buf, ' ')
		}
	}
    // file and line number (if corresponding flags are provided)
	if l.flag&(Lshortfile|Llongfile) != 0 {
		if l.flag&Lshortfile != 0 {
			short := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
			file = short
		}
		*buf = append(*buf, file...)
		*buf = append(*buf, ':')
		itoa(buf, line, -1)
		*buf = append(*buf, ": "...)
	}
    // l.prefix (if it's not blank and Lmsgprefix is set)
	if l.flag&Lmsgprefix != 0 {
		*buf = append(*buf, l.prefix...)
	}
}
~~~

最后就是一些工具类：**为什么要有这个工具类？**

~~~go
// Cheap integer to fixed-width decimal ASCII. Give a negative width to avoid zero-padding.
func itoa(buf *[]byte, i int, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}
~~~

比如要输出年月日信息，`year, month, day := t.Date()` 得到的都是 int 类型的值，但是装载输出内容的容器是 `[]byte`，而且最终输出的都是按照 string 内容输出。也就是**最终将 int 值转化成对应的 ASCII 值**。



疑惑：

1. 为什么 `Release lock while getting caller info - it's expensive.`
