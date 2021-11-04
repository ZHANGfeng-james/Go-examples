package v1

import (
	"fmt"
	"reflect"
)

func changeValue(v *int) {
	value := reflect.ValueOf(v) // *int --> reflect.Value
	ele := value.Elem()         // *int --> int

	if ele.CanSet() {
		ele.SetInt(20) // SetInt 设置 int 类型值
	}
}

type person struct {
}

func createNewIntValue() {
	var v int = 10

	// 创建一个 *int --> reflect.Value，获取值并修改值
	value := reflect.New(reflect.TypeOf(v)) // *int --> reflect.Value
	instance := value.Interface()           // reflect.Value --> interface{}值

	person := person{}
	ptr := &person
	fmt.Println(reflect.Indirect(reflect.ValueOf(ptr)).Type().Name())

	//FIXME Body 部分的 req.argv 是如何被解析出来的？RPC 中 encoding/gob 的入参是 interface{}

	var tmp int = 20
	instance = &tmp

	fmt.Println(instance, value.Elem()) // value.Elem() 获取 Pointer 指向的值
}
