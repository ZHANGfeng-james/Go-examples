package reslice

import "fmt"

type T0 struct {
	x int
}

func (*T0) M0() {

}

type T1 struct {
	y int
}

func (T1) M1() {

}

type T2 struct {
	z int
	T1
	*T0
}

func (*T2) M2() {

}

type Q *T2

func Init() {

}

func Test() {
	fmt.Println("Running!")
}
