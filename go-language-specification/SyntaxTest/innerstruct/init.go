package innerstruct

import "fmt"

func InitTest1() {
	type A struct {
		a int
		b int
		c string
	}

	value := A{
		1, 2, "",
	}
	fmt.Println(value)

}
