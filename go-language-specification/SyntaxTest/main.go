package main

import "fmt"

func main() {
	value := []byte("\xc5")
	fmt.Println(len(value))
}
