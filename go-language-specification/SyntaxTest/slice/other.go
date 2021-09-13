package reslice

import "fmt"

const PI = 3.14

var value = initValue()

func init() {
	fmt.Println("init")
}

func initValue() int {
	fmt.Println("reslice package initValue func...")
	return 1
}

func Test() {
	fmt.Println("Test")
}
