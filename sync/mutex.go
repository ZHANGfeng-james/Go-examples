package sync

import (
	"log"
	"sync"
)

type Counter struct {
	sync.Mutex
	Count int
}

func copyMutexInstance() {
	counter := Counter{}
	counter.Lock()
	defer counter.Unlock()
	counter.Count++
	foo(&counter)
}

func foo(lock sync.Locker) {
	lock.Lock()
	defer lock.Unlock()
	log.Println("foo")
}
