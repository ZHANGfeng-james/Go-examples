package list

import (
	"container/list"
	"fmt"
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

func EmptyList() {
	var list *list.List // list 的初始值是 nil

	list.Init()

	list.PushBack("a")
	fmt.Println(list.Len())
}

func MarkIsNil() {
	list := list.New()

	aEle := list.PushBack("a")
	list.InsertAfter("b", aEle)

	list.InsertAfter("c", nil)

	fmt.Println("len:", list.Len())
}
