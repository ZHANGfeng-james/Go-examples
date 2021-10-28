package builtin

import "log"

func newUsage() {
	var i *int
	i = new(int)
	log.Printf("value:%d", *i)

	defer func() {
		if err := recover(); err != nil {
			// runtime error: invalid memory address or nil pointer dereference
			log.Println(err)
		}
	}()

	var ptr *int
	log.Printf("value of ptr:%d", *ptr)
}

func newChannel() {
	ch := new(chan int)
	log.Printf("ch's type:%T", ch)

	kVal := new(map[string]int)
	log.Printf("kVal's type:%T", kVal)
	defer func() {
		if err := recover(); err != nil {
			// assignment to entry in nil map
			log.Println(err)
		}
	}()

	(*kVal)["1"] = 1
}
