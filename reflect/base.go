package v1

import (
	"fmt"
	"reflect"
)

type Account struct {
	Username string `geekorm:"PRIMARY KEY"`
	Age      int8
}

func (account *Account) GetAge() {

}

func (account Account) GetUsername() {

}

func main() {
	ptrTest()
	normalTest()
	fmt.Println("<--- ptrIndirect --->")
	ptrIndirect()
	fmt.Println("<-reflectNew->")
	reflectNew()
	fmt.Println("<--compareTypeOf-->")
	compareTypeOf()
	fmt.Println("<--recordValues-->")
	recordValues()
	fmt.Println("<--elemAndInterface-->")
	elemAndInterface()
	fmt.Println("<--indirect--->")
	indirect()
}

func ptrTest() {
	ptr := &Account{}
	fmt.Printf("%T\n", ptr) // *main.Account

	typ := reflect.TypeOf(ptr)
	fmt.Println(typ.Name(), typ.NumMethod(), typ.Kind()) // 空

	value := reflect.ValueOf(ptr)
	fmt.Println(value.Kind(), value.Type(), value.NumMethod())
}

func normalTest() {
	obj := Account{}
	fmt.Printf("%T\n", obj)

	value := reflect.ValueOf(obj)
	fmt.Println(value.Kind(), value.Type().Name(), value.NumField(), value.NumMethod())

	typ := value.Type()
	fmt.Println(typ.Name(), typ.Kind(), typ.NumField(), typ.NumMethod())
}

func ptrIndirect() {
	ptr := &Account{
		Username: "Katyusha",
		Age:      18,
	}
	typ := reflect.Indirect(reflect.ValueOf(ptr)).Type()
	fmt.Println(typ.Name(), typ.NumField())

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fmt.Println(field.Name, field.Tag)
	}
}

func reflectNew() {
	ptr := Account{}

	// start: reflect.Type
	typ := reflect.TypeOf(ptr)

	valuePtr := reflect.New(typ)

	// end: reflect.Value
	value := reflect.Indirect(valuePtr)
	fmt.Println(value.Kind(), value.NumField(), value.NumMethod())
}

type Password struct {
	life    int
	content string
}

func compareTypeOf() {
	ptr := &Account{}
	fmt.Printf("%T\n", ptr) // *main.Account

	typ := reflect.TypeOf(ptr)

	fmt.Println(typ == reflect.TypeOf(&Account{}))
}

func recordValues() {
	account := &Account{
		Username: "Katyusha",
		Age:      18,
	}

	// 获取 Fields 数组
	typ := reflect.Indirect(reflect.ValueOf(&Account{})).Type()
	fmt.Println(typ.NumField())

	value := reflect.Indirect(reflect.ValueOf(account))
	for i := 0; i < typ.NumField(); i++ {
		// 依据 field.Name 获取对应的值，此次 Field 必须是可导出的
		v := value.FieldByName(typ.Field(i).Name).Interface()
		fmt.Printf("fieldName[%d]=%s, value:%v\n", i, typ.Field(i).Name, v)
	}
}

func elemAndInterface() {
	var accounts []Account

	var ptr interface{}
	ptr = &accounts

	destSlice := reflect.Indirect(reflect.ValueOf(ptr)) // reflect.Value []Acccount
	destType := destSlice.Type().Elem()                 // Account
	fmt.Println(destType.Name())

	value := reflect.New(destType).Elem() // reflect.Value
	value.Interface()                     // interface{}
}

func indirect() {
	account := Account{
		Username: "Katyusha",
		Age:      18,
	}
	value := reflect.Indirect(reflect.ValueOf(account))
	fmt.Println(value.FieldByName("Username"))

	fmt.Println(value.Type().Name())
}
