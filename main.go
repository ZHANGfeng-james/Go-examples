package main

import "log"

type Person struct {
	name string
}

func (p *Person) getName() string {
	return p.name
}

func (p *Person) setName(name string) {
	p.name = name
}

func main() {
	p := Person{name: "Katyusha"}
	log.Printf("name:%s", p.getName())

	p.setName("Authur")
	log.Printf("name:%s", p.getName())
}
