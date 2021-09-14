`time.go` 的源文件包括：Time、Duration 类型，下面具体分析其含义！

最基本的，和时间单位秒的**换算关系**：

~~~go
1s = 1000 ms = 1000,000 us = 1000, 000, 000 ns
~~~

在典型 PC 机上各种操作的**近似时间（相对值）**：

~~~go
执行典型指令     　　　　　　　　　　  1/1,000,000,000 秒 =1 纳秒
从一级缓存中读取数据 　　　 　　　　   0.5 纳秒
分支预测错误 　　　　　　    　　　　  5 纳秒
从二级缓存中读取数据 　　　　　　　    7 纳秒

互斥锁定 / 解锁 　　　　　　 　　　　  25 纳秒

从主存储器中读取数据 　　    　　　　  100 纳秒 

在 1Gbps 的网络中发送 2KB 数据 　　   20,000 纳秒

从内存中读取 1MB 数据 　　　　　　     250,000 纳秒
从新的磁盘位置读取数据 ( 寻轨 ) 　　   8,000,000 纳秒
从磁盘中读取 1MB 数据 　　　　　　     20,000,000 纳秒

在美国向欧洲发包并返回 　　　　　　     150 毫秒 =150,000,000 纳秒
~~~

time 包提供了测量（measuring not setting，用于测量目的，不受时间修改的影响）和显示（displaying）时间的方法。**前提和原则**：日历历法计算始终采用**公历**，且没有**闰秒**。

**在 Linux 世界**里有 4 种**时钟类型**：

1. `CLOCK_REALTIME`：就是非常出名的 `Wall CLock`，就是实际的时间（20xx 年 xx 月 xx 日 xx 时 xx 分 xx.xx 秒），也就是挂钟时间。`CLOCK_REALTIME` 是**可以设置的**，用户可以使用 date 命令或是系统调用去修改。当系统休眠（Suspend）时，`CLOCK_REALTIME` 仍然会运行正常，但系统恢复时，Kernel 去做补偿。
2. `CLOCK_MONOTONIC`：即单调时间，从**某个时间点**开始到现在已经流逝的时间。用户不能修改这个时间，但是当系统进入休眠时，`CLOCK_MONOTONIC` 是不会增加的。`CLOCK_MONOTONIC` 相当于是一个计时器，但系统开机启动的时候重新开始计时，但休眠时暂停计时；当系统关机时，停止计时。
3. `CLOCK_MONOTONIC_RAW`：
4. `CLOCK_BOOTTIME`：

对应的，**在 Go 的 time 包中**有 2 个概念：

* Wall Clock：就是一般意义上的时间，就像墙上挂钟所指示的时间；
* Monotonic Clock：字面意思是**单调时间**，是指从某个点开始后流逝的时间，比如系统启动后。

为了不去拆分 API，Go 标准库中的使用 Time 实例包含了上述两种时间信息，即：`a wall clock reading` 和 `a monotonic clock reading`，可使用 `time.Now` 返回。比如下述代码：

~~~go
func main() {
	ch := make(chan struct{})

	start := time.Now()
	time.Sleep(5 * time.Second)
	elapsed := time.Now().Sub(start)
    // 5.0007701s
	fmt.Println(elapsed)

	sendMessage(ch)
	<-ch
}
~~~

`time.Sleep` 用户模拟计算机操作耗时，即使在定时操作期间修改了 Wall Clock，上述结果输出也始终是大约 5 秒的时间间隔。time 包中另外的一些方法，比如 `time.Since(start)`、`time.Until(deadline)`、`time.Now().Before()` 同样不会受到 Wall Clock 改变的影响。

`time.Now` 返回的 Time 实例包含了 `Monotonic Clock` 读数，`t.Add` 获得结果将会让 `the wall clock readings` 和 `the monotonic clock readings` 增加相同的间隔时间：

~~~go
func main() {
	ch := make(chan struct{})

	start := time.Now()
    // 2021.06.20 14:25:48
	fmt.Println(start.Format("2006.01.02 15:04:05"))
    // 2021.06.20 14:25:53
	fmt.Println(start.Add(5 * time.Second).Format("2006.01.02 15:04:05"))

	sendMessage(ch)
	<-ch
}
~~~

`t.AddDate(y, m, d)`，`t.Round(d)`，`t.Truncate(d)` 这些方法在计算时使用的是 `the wall clock reading`，因此在计算时，会**剔除**掉 `the monotonic clock reading`。相同的，`t.In`、`t.Local` 和 `t.UTC` 只是用于解释 `the wall clock reading`，因此，它们会从结果中剔除掉 `the monotonic clock reading`。剔除掉 `the monotonic clock reading` 标准的方式是调用 `t.Round(0)`。

如果 Time 实例 t 和 u 同时都包含了 `the monotonic clock reading`，那么 `t.After(u)`、`t.Before(u)`、`t.Equal(u)` 和 `t.Sub(u)` 在计算时只会使用该读数，而忽略 `the wall clock reading`。但只要 t 和 u 中的任何一个实例不包含 `the monotonic clock reading`，那么上面这些方法就会使用 `the wall clock reading` 而忽略另一个实例的 `the monotonic clock reading`。

在一些操作系统中，如果系统进入休眠时，`the monotonic clock reading` 会停止计数，直到系统被唤醒时才会继续增加。在这样的系统中，`t.Sub(u)` 不会准确地反映出 t 和 u 之间的实际的时间间隔。

因为在**当前进程之外（也就是说，这个单调时间是和当前进程相关的）**使用 `the monotonic clock reading`  是没有任何意义的，因此序列化的方法 `t.GobEncode`、`t.MarshalBinary`、`t.MarshalJSON` 和 `t.MarshalText` 会剔除掉 `the monotonic clock reading` 数值。同样的，`t.Format` 也不为该数值提供任何格式化信息。相类似的原因，`time.Date`、`time.Parse`，`time.ParseInLocation` 和 `time.Unix`，以及一些反序列化方法比如 `t.GobDecode` ，`t.UnmarshalBinary`、`t.UnmarshalJSON` 和 `t.UnmarshalText` 总是创建没有 `the monotonic clock reading` 的时间。

在 Go 语言中的 `==` 运算符作用在 Time 实例上时，不仅比较 `the wall clock reading` ，还比较 Location 和 `the monotonic clock reading`。

~~~go
func main() {
	ch := make(chan struct{})

	start := time.Now()

	// 2021-06-20 15:01:45.4863014 +0800 CST m=+0.003997201
	fmt.Println(start.String())
	// 2021-06-20 15:01:45.4863014 +0800 CST
	fmt.Println(start.Round(0).String())

	// true
	fmt.Println(start == start)
	// false
	fmt.Println(start == start.Round(0))

	sendMessage(ch)
	<-ch
}
~~~

`start.String()` 的输出结果中，`m=+0.003997201` 就是 `the monotonic clock reading`。

# 1 Time

一个 Time 实例，表示的是按照**纳秒**为精度的时间值。

App 中使用 Time 时，应该直接按照 Time 类型值存储时间，而不是使用其指针 `*time.Time`。也就是说，Time 类型的变量，或者作为结构体的字段时，都应该是 `time.Time` 而不应该是 `*time.Time`。

> 对 `time.Time` 类型使用时，做了严格的限制，为什么？为什么不能使用 `*time.Time` 类型？

一个 Time 实例可以在多个 goroutine 中**并发**使用，但是在 `GobDecode`、`UnmarshalBinary`、`UnmarshalJSON` 和 `UnmarshalText` 却是**非并发安全**的。

Time 实例可以使用 Before、After、Equal 方法做比较。Sub 方法，其结果返回的是一个 Duration 类型实例，Add 方法作用在 Time 和 Duration 实例上，其结果返回的是一个 Time 实例。

Time 的**零值时间**代表的是：`January 1, year 1, 00:00:00.000000000 UTC`。实际上，这个零值时间并不具备有实际的含义，`IsZero` 方法用一种简单的方法判断 Time 实例是否显示初始化：

~~~go
package main

import (
	"fmt"
	"time"
)

func main() {
	var zeroTime time.Time
	fmt.Println(zeroTime.IsZero())

	fmt.Println(zeroTime)
}
true
0001-01-01 00:00:00 +0000 UTC
~~~

每一个 Time 实例都附带有一个 Location 信息，但在调用 Format、Hour 和 Year 方法时，会使用到该 Location 信息。方法 Local、UTC 和 In 都会返回一个包含有指定 Location 的 Time 实例。改变 Time 中的 Location 信息，仅仅会改变 Time 的显示，并不会改变 Time 代表的时间，也不会对早先的时间计算造成影响。

> 协调世界时，又称世界统一时间、世界标准时间、国际协调时间。由于英文（CUT）和法文（TUC）的缩写不同，作为妥协，简称UTC。
>
> 世界统一时间，不属于任意时区！**时区**(Time Zone)是地球上的**区域使用同一个时间定义**。1884年在华盛顿召开国际经度会议时，为了克服时间上的混乱，规定将全球划分为24个时区。在中国采用首都北京所在地东八区的时间为全国统一使用时间。

使用 `GobEncode`、`MarshalBinary`、`MarshalJSON` 和 `MarshalText` 存储一个表征 Time 的值时，会存储 `Time.Location` 的偏移量，而并不是 location 的名称，因此这种方式会都是 `DST` 信息。

> 夏令时，表示**为了节约能源**，人为规定时间的意思。也叫夏时制，夏时令（`Daylight Saving Time：DST`），又称“日光节约时制”和“夏令时间”。

需要指出的是 Go 中的 `==` 操作符作用在 Time 实例上时，不仅仅比较的是 Time 实例，还要比较 Location 和 `the monotonic clock reading`。因此，如果没有办法保证所有 Time 实例设置了相同的 Location，Time 实例不能用作 Map 和数据库的 key 值。对于 Time 实例的比较，更加合适的做法是使用 `Equal()` 方法而不是 `==` 操作符，只有当其中一个参数具有 `monotonic clock` 读数时才能正确处理。

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

wall 和 `ext` 属性值，包含了以秒为单位的 `the wall time` 值，以及以纳秒为单位的 `the wall time` 值，以及可选的以纳秒为单位的 `the monotonic time` 值。3 个属性是：

| 属性  |   类型    | 说明                                                         |
| :---: | :-------: | :----------------------------------------------------------- |
| wall  |  uint64   | 从高比特位到低位，包含 1 个 bit 的标志位，表示是否包含 monotonic 值；33-bit 秒值；30-bit 纳秒值，范围是[0, 999999999] |
| `ext` |   int64   | 标志位值是 0，上述 33-bit 秒值必须是 0，且 `ext` 值存储的是从零值来时的秒值；标志位是 1，上述 33-bit 秒值代表的是从 1885 年 1 月 1 日开始的值；`ext` 表示**从进程开始**的纳秒值 |
| `loc` | *Location | 时区 ，默认为 `nil` ，也就是 UTC ，或者说格林威治时间        |

也就是说，`loc` 指定了用于确定于此时间对应的分钟、小时、月份、日期和年份的时区。













