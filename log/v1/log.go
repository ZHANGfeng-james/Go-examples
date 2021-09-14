package main

import (
	"log"
	"os"
)

func main() {
	logger := log.New(os.Stdout, "---", log.LstdFlags)
	logger.Print("Hello!")
}
