# 1 Time

Go time 标准库中定义的 Time 类型，包含如下信息：

~~~go
type Time struct {
	// wall and ext encode the wall time seconds, wall time nanoseconds,
	// and optional monotonic clock reading in nanoseconds.
	//
	// From high to low bit position, wall encodes a 1-bit flag (hasMonotonic),
	// a 33-bit seconds field, and a 30-bit wall time nanoseconds field.
	// The nanoseconds field is in the range [0, 999999999].
	// If the hasMonotonic bit is 0, then the 33-bit field must be zero
	// and the full signed 64-bit wall seconds since Jan 1 year 1 is stored in ext.
	// If the hasMonotonic bit is 1, then the 33-bit field holds a 33-bit
	// unsigned wall seconds since Jan 1 year 1885, and ext holds a
	// signed 64-bit monotonic clock reading, nanoseconds since process start.
	wall uint64
	ext  int64

	// loc specifies the Location that should be used to
	// determine the minute, hour, month, day, and year
	// that correspond to this Time.
	// The nil location means UTC.
	// All UTC times are represented with loc==nil, never loc==&utcLoc.
	loc *Location
}
~~~

就是说，一个 Time 实例，默认包含了 Location 实例，表示的是**特定时区**的时间信息。也只有在描述时间时，附带上时区（地区信息）才是有效的。

获取 Time 实例的方式有以下几种：

~~~go
// Now returns the current local time.
func Now() Time {
	sec, nsec, mono := now()
	mono -= startNano
	sec += unixToInternal - minWall
	if uint64(sec)>>33 != 0 {
		return Time{uint64(nsec), sec + minWall, Local}
	}
	return Time{hasMonotonic | uint64(sec)<<nsecShift | uint64(nsec), mono, Local}
}

func Parse(layout, value string) (Time, error) {
	return parse(layout, value, UTC, Local)
}
~~~











获取**当前时间**，并转化为**指定格式**：

~~~go
package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func main() {
    formatTime := time.Now().Format("2006-01-02_15-04-05")
	file, err := os.Create("Snipaste_" + formatTime + ".png")
}
~~~

最终获取的时间字符串：`Snipaste_2021-07-22_11-58-41.png`









~~~go
// Nanosecond returns the nanosecond offset within the second specified by t,
// in the range [0, 999999999].
func (t Time) Nanosecond() int {
	return int(t.nsec())
}

// UnixNano returns t as a Unix time, the number of nanoseconds elapsed
// since January 1, 1970 UTC. The result is undefined if the Unix time
// in nanoseconds cannot be represented by an int64 (a date before the year
// 1678 or after 2262). Note that this means the result of calling UnixNano
// on the zero Time is undefined. The result does not depend on the
// location associated with t.
func (t Time) UnixNano() int64 {
	return (t.unixSec())*1e9 + int64(t.nsec())
}
~~~

如果单单从方法的返回结果来看，确实是有不同的：int 和 int64 的类型。

通过如下的示例代码验证：

~~~go
package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan struct{})

	start := time.Now().Nanosecond()
	// 5995900
	time.Sleep(time.Millisecond * 5)
	end := time.Now().Nanosecond()

	fmt.Println(end - start)

	sendMessage(ch)
	<-ch
}

func sendMessage(ch chan<- struct{}) {
	go func() {
		ch <- struct{}{}
	}()
}

~~~





# 2 Duration

和时间对象相关的，还有 Duration 类型，表示持续的时长（start-end 之间的时间间隔）。

~~~go
// A Duration represents the elapsed time between two instants
// as an int64 nanosecond count. The representation limits the
// largest representable duration to approximately 290 years.
type Duration int64
~~~

获取 Duration 实例：

~~~go
// ParseDuration parses a duration string.
// A duration string is a possibly signed sequence of
// decimal numbers, each with optional fraction and a unit suffix,
// such as "300ms", "-1.5h" or "2h45m".
// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
func ParseDuration(s string) (Duration, error) {
	// [-+]?([0-9]*(\.[0-9]*)?[a-z]+)+
	orig := s
    ...
}
~~~



# 3 Location

Location 类型表示的是和 Time 相关的时区信息。

~~~go
func LoadLocation(name string) (*Location, error)
~~~

如果给定的 name 是空的，或者是 `UTC`，上述函数将返回 UTC 时区实例；如果给定的是 `Local`，则返回操作系统当前时区。比如：

~~~go
func main() {
	location, _ := time.LoadLocation("Asia/Shanghai")

	inputTime := "2029-09-04 12:02:33"
	layout := "2006-01-02 15:04:05"
	t, _ := time.Parse(layout, inputTime)
	// t, _ := time.ParseInLocation(layout, inputTime, location)

	dateTime := time.Unix(t.Unix(), 0).In(location).Format(layout)

    // 输入时间：2029-09-04 12:02:33, 输出时间:2029-09-04 20:02:33
    // 输入时间：2029-09-04 12:02:33, 输出时间:2029-09-04 12:02:33
	fmt.Printf("输入时间：%s, 输出时间:%s\n", inputTime, dateTime)
}
~~~

`Asia/Shanghai` 就是亚洲上海的时区实例。





