package main

import (
	"fmt"
	"reflect"
)

type Account struct {
	username string
	age      int8
}

func main() {
	ptr := &Account{}
	typ := reflect.Indirect(reflect.ValueOf(ptr)).Type()
	fmt.Println(typ.Name())

	obj := Account{}
	value := reflect.ValueOf(obj)
	typ = value.Type()
	fmt.Println(typ.Name())
}
