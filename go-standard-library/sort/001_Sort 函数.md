> sort 包提供了用于对切片和用户自定义的集合进行**排序**的原语。



~~~go
func Search(n int, f func(int) bool) int
~~~

Search 使用二分查找方法，在 **[0, n)** 中找到**最小的**、能够满足 f(i) 值为 true 的值，其前提是：**f(i) == true**，意味着 f(i+1) == true。需要注意的是：**其返回值 i 的取值范围是 `[0, n)`！**如果首次取值后，f(i) 返回的值是 false，则继续取 i + 1（若此时该值 < n 时）。

如果所有的取值都不能让 f(i) 的结果为 true，即不能满足 f(i) 的条件，返回值则为 n。

Search 函数一般使用在排序集合中的搜索问题上，此时参数 f 就是一个**闭包**，其中包含了待搜索的内容集合，而且内容集合是有序的。

比如：给定的切片是**升序的**，调用 `Search(len(data), func(i int) bool { return data[i]>=23})` 的结果返回的是能够满足 `data[i]>=23` 的最小位置索引。如果待搜索内容集合是**降序的**，则需要使用的是 `<=` 操作符。

比如下面是一个**升序集合**的例子：

~~~go
x := 23
i := sort.Search(len(data), func(i int) bool { return data[i] >= x })
if i < len(data) && data[i] == x {
	// x is present at data[i]
} else {
	// x is not present in data,
	// but i is the index where it would be inserted.
}
~~~

官方给出了一个二分法猜数值的例子：

~~~go
func GuessingGame() {
	var s string
	fmt.Printf("Pick an integer from 0 to 100.\n")
	answer := sort.Search(100, func(i int) bool {
		fmt.Printf("Is your number <= %d? ", i)
		fmt.Scanf("%s", &s)
		return s != "" && s[0] == 'y'
	})
	fmt.Printf("Your number is %d.\n", answer)
}
~~~

但实际上是有问题的！问题出在第6行的 Scanf 函数中。比如首次键入了 `yCRLF`，也就是一个英文字符 y 和一个回车换行符，但是 fmt.Scanf 在第二次循环读取时，读取的是回车换行符。这是不符合输入要求的，其解决办法就是过滤回车换行符，从回车换行符下一个字符开始读取：

~~~go
func GuessingGame() {
	var s string
	fmt.Printf("Pick an integer from 0 to 100.\n")
	answer := sort.Search(100, func(i int) bool {
		fmt.Printf("Is your number <= %d? ", i)

		stdin := bufio.NewReader(os.Stdin)
		fmt.Fscan(stdin, &s)
		stdin.ReadString('\n')
		return s != "" && s[0] == 'y'
	})
	fmt.Printf("Your number is %d.\n", answer)
}
~~~

运行结果如下：

~~~go
PS E:\go_developer_roadmap\ProgrammingLanguage\Go Standard Interface\GoUsage> go run main.go
Pick an integer from 0 to 100.
Is your number <= 50? y
Is your number <= 25? n
Is your number <= 38? y
Is your number <= 32? y
Is your number <= 29? y
Is your number <= 27? n
Is your number <= 28? y
Your number is 28.
~~~

再比如：

~~~go
package main

import (
	"fmt"
	"sort"
)

func main() {
    // []int 内容是升序排列的
	a := []int{1, 3, 6, 10, 15, 21, 28, 36, 45, 55}
	x := 6

	i := sort.Search(len(a), func(i int) bool { return a[i] >= x })
	if i < len(a) && a[i] == x {
		fmt.Printf("found %d at index %d in %v\n", x, i, a)
	} else {
		fmt.Printf("%d not found in %v\n", x, a)
	}
}
~~~

再比如一个降序的例子：

~~~go
package main

import (
	"fmt"
	"sort"
)

func main() {
	a := []int{55, 45, 36, 28, 21, 15, 10, 6, 3, 1}
	x := 6

	i := sort.Search(len(a), func(i int) bool { return a[i] <= x })
	if i < len(a) && a[i] == x {
		fmt.Printf("found %d at index %d in %v\n", x, i, a)
	} else {
		fmt.Printf("%d not found in %v\n", x, a)
	}
}
~~~

# 1 源代码示例

net/http/server.go 中：

~~~go
// Handle registers the handler for the given pattern.
// If a handler already exists for pattern, Handle panics.
func (mux *ServeMux) Handle(pattern string, handler Handler) {
	mux.mu.Lock()
	defer mux.mu.Unlock()

	if pattern == "" {
		panic("http: invalid pattern")
	}
	if handler == nil {
		panic("http: nil handler")
	}
	if _, exist := mux.m[pattern]; exist {
		panic("http: multiple registrations for " + pattern)
	}

	if mux.m == nil {
		mux.m = make(map[string]muxEntry)
	}
	e := muxEntry{h: handler, pattern: pattern}
	mux.m[pattern] = e
	if pattern[len(pattern)-1] == '/' {
        // mux.es 类型 []muxEntry（slice of entries sorted from longest to shortest.）
		mux.es = appendSorted(mux.es, e)
	}

	if pattern[0] != '/' {
		mux.hosts = true
	}
}

func appendSorted(es []muxEntry, e muxEntry) []muxEntry {
    // es 类似于查找容器，可以为空（若为空，则 n 值为 0，i 的值也肯定就为 0）
	n := len(es)
	i := sort.Search(n, func(i int) bool {
        // 从使用的操作符 < 就可以看出，容器是降序排列的
		return len(es[i].pattern) < len(e.pattern)
	})
	if i == n {
		return append(es, e)
	}
	// we now know that i points at where we want to insert
	es = append(es, muxEntry{}) // try to grow the slice in place, any entry works.
	copy(es[i+1:], es[i:])      // Move shorter entries down
	es[i] = e
	return es
}
~~~

上述 appendSorted 可以看成是一个容器的创建过程，而且是**降序排列**的。特别是，appendSorted 函数使用 sort.Search 用于查找待插入的位置索引。

