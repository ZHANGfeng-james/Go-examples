package reslice

import (
	"fmt"

	pkgname "opensource.com/syntaxtest/pkg"
)

func InitTest1() {
	var slice1 []int            // slice is nil
	var slice2 = make([]int, 0) // slice is not nil!

	if slice1 == nil {
		fmt.Println("slice1 is nil")
	}

	if slice2 == nil {
		fmt.Println("slice2 is nil")
	}

	p1 := &[]int{}   // p1 points to an initialized, empty slice with value []int{} and length 0
	p2 := new([]int) // p2 points to an uninitialized slice with value nil and length 0
	fmt.Println(*p1, *p2)
	if *p1 == nil {
		fmt.Println("*p1 is nil")
	}

	if *p2 == nil {
		fmt.Println("*p2 is nil")
	}
}

func InitTest2() {
	tmp := [...]int{1, 2, 3, 4}
	// tmp[:2] 切片操作，而不是数组操作
	fmt.Printf("%T.\n", tmp[:2])
}

func InitTest3() {
	value := 3
	if value == ([]int{1, 2, 3, 4}[2]) {
		fmt.Println("value:", value)
	}
}

func InitTest4() {
	vowels := [128]bool{'a': true, 'e': true, 'i': true, 'o': true, 'u': true, 'y': true}
	fmt.Println(vowels['b'])
}

func InitTest5() {
	pkgname.TestInit()
}
