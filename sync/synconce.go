package sync

import (
	"log"
	"sync"
)

func onceUsage() {
	var o sync.Once

	o.Do(func() {
		log.Println("Hello")
	})

	o.Do(func() {
		log.Println("Hello Again!")
	})

	var other sync.Once
	other.Do(func() {
		log.Println("Hello Again!")
	})
}

func onceCopy() {
	var once sync.Once

	once.Do(func() {
		log.Println("ready to copy!")
	})

	twoOnce := once
	twoOnce.Do(func() {
		log.Println("copy success")
	})
}
