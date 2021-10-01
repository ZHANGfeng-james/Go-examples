Go 语言标准库 container/list 实现了一个**双向链表**，类似于：

![](./doubly_linked_list.png)

在这个双向链表中，**每个节点都包含了 prev 和 next 字段**，分别指向了**前置节点**和**后续节点（接下来的节点）**。自然的，如果是 list 的首节点，prev 的值是 nil；如果是末尾节点，next 的值是 nil。

在 container/list 包中，**结构体类型**：

~~~go
ant@MacBook-Pro Go-examples-with-tests % go doc container/list |grep "^type"
type Element struct{ ... }
type List struct{ ... }
~~~

如果想要遍历 list 的各个节点，按照如下方式：

~~~go
...
for node := list.Front(); node != nil; node = node.Next() {
    // do something with node.Value
}
...
~~~

container/list 的使用示例代码：

~~~go
package list

import (
	"container/list"
	"strings"
)

func UseList() interface{} {
	list := list.New()

	aNode := list.PushFront("a")
	dNode := list.PushBack("d")

	list.InsertAfter("b", aNode)
	list.InsertBefore("c", dNode)

	var result strings.Builder
	for ele := list.Front(); ele != nil; ele = ele.Next() {
		result.WriteString(ele.Value.(string))
	}
	return result.String()
}
~~~

# 1 Element

container/list 包中，封装**双向链表节点**的数据结构是 Element（不是 Node）：

~~~go
// Element is an element of a linked list.
type Element struct {
	// Next and previous pointers in the doubly-linked list of elements.
	// To simplify the implementation, internally a list l is implemented
	// as a ring, such that &l.root is both the next element of the last
	// list element (l.Back()) and the previous element of the first list
	// element (l.Front()).
	next, prev *Element

	// The list to which this element belongs.
	list *List

	// The value stored with this element.
	Value interface{}
}
~~~

Value 是可导出的，也就是说 ele.Value 就是节点中存放的内容（值）。

Element 节点相关的 2 个方法，分别是：

* `func (e *Element) Next() *Element`：也就是线性结构中，当前节点的下一个节点，或者是 nil；
* `func (e *Element) Prev() *Element`：当前节点的上一个节点，或者是 nil。

# 2 List

另外一个重要的数据结构是：

~~~go
// List represents a doubly linked list.
// The zero value for List is an empty list ready to use.
type List struct {
	root Element // sentinel list element, only &root, root.prev, and root.next are used
	len  int     // current list length excluding (this) sentinel element
}
~~~

为了构造出一个 List，与之关联的**全局方法**是：

~~~go
func New() *List
~~~

该方法用于构造出一个 List，返回值类型是 `*list.List`，相当于是对 List 的**初始化操作**。下面这种用法是错误的：

~~~go
func EmptyList() {
	var list *list.List // list 的初始值是 nil

	list.Init() // runtime error: invalid memory address or nil pointer dereference

	list.PushBack("a")
	fmt.Println(list.Len())
}
~~~

Go 语言中指针类型变量的初始值是 nil，在第 4 行会抛出空指针异常。因此，在使用前，必须调用 `list.New()` 获取一个 `*list.List` 变量，这样才代表**一个可用的双端链表**。

下面依次说明和 `*list.List` 相关的方法：

* `func (l *List) Back() *Element`：返回链表的末尾节点；
* `func (l *List) Front() *Element`：返回链表的首节点；
* `func (l *List) Init() *List`：用于初始化链表（New 函数调用了），或者是清空链表；
* `func (l *List) InsertAfter(v interface{}, mark *Element) *Element`：mark 不能是 nil，否则会抛出空指针异常。该方法的作用是在 mark 之后插入 Value 是 v 的 Element（新创建节点）。如果 mark 不属于当前操作的 list，则不会有任何改变。
* `func (l *List) InsertBefore(v interface{}, mark *Element) *Element`
* `func (l *List) Len() int`：返回 list 的节点个数，时间复杂度是 O(1)
* `func (l *List) MoveAfter(e, mark *Element)`：e 和 mark 不能是 nil，如果 e 或者 mark 不属于当前操作的 list，将不会有任何改变。该方法的效用是：将 e 移到 mark 节点之后。
* `func (l *List) MoveBefore(e, mark *Element)`
* `func (l *List) MoveToBack(e *Element)`
* `func (l *List) MoveToFront(e *Element)`
* `func (l *List) PushBack(v interface{}) *Element`
* `func (l *List) PushBackList(other *List)`：在 list 的末尾新加入 other 的一份拷贝，相当于是在 list 的末尾另外接入了 other 双端链表。即 `list.Len() += other.Len()`，other 不能为 nil！
* `func (l *List) PushFront(v interface{}) *Element`
* `func (l *List) PushFrontList(other *List)`
* `func (l *List) Remove(e *Element) interface{}`：在 list 移除 e 这个节点，如果 e 是不属于 list 的，将不会有任何变化。

关于上述方法，一个很有意思的结论是：

~~~go
// Remove removes e from l if e is an element of list l.
// It returns the element value e.Value.
// The element must not be nil.
func (l *List) Remove(e *Element) interface{} {
	if e.list == l {
		// if e.list == l, l must have been initialized when e was inserted
		// in l or l == nil (e is a zero Element) and l.remove will crash
		l.remove(e)
	}
	return e.Value
}
~~~

在方法中根本没有判断 e 是否是 nil，此时，如果 e 是 nil，则会直接 crash！

**上述方法中，绝大多数都要求入参不能是 nil**。

