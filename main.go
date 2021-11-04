package main

import (
	"log"
	"sync"
)

func main() {
	var kVal sync.Map

	kVal.Store("1", 1)
	kVal.Range(func(key, value interface{}) bool {
		log.Printf("%s --> %v", key, value)
		return true
	})

	kVal.Delete("1")

}
