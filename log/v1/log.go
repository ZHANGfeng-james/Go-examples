package main

import (
	"fmt"
	"log"
	"os"
)

func prefixUsage() {
	var helloStr = "Hello, log!"
	stdFlags := log.Flags() // LstdFlags = Ldate | Ltime
	log.Print(helloStr, " stdFlags:", stdFlags)

	log.SetPrefix("---")

	var prefixFlag map[int]string = make(map[int]string)
	prefixFlag[log.LUTC] = "log.LUTC"
	prefixFlag[log.Ldate] = "log.Ldate"
	prefixFlag[log.Ltime] = "log.Ltime"
	prefixFlag[log.Lmicroseconds] = "log.Lmicroseconds"
	prefixFlag[log.LstdFlags] = "log.LstdFlags"
	prefixFlag[log.Llongfile] = "log.Llongfile"
	prefixFlag[log.Lshortfile] = "log.Lshortfile"
	prefixFlag[log.Lmsgprefix] = "log.Lmsgprefix"

	for prefix, str := range prefixFlag {
		log.SetFlags(prefix)
		fmt.Printf("%+20s:", str)
		log.Print(helloStr)
	}
}

// usePrefixCompose
func usePrefixCompose() {
	logger := log.New(os.Stdout, "---", log.LstdFlags)
	logger.Print("Hello!") // ---2021/10/15 21:18:22 Hello!

	logger = log.New(os.Stdout, "---", log.LstdFlags|log.Lmsgprefix)
	logger.Print("Hello!") // 2021/10/15 21:18:22 ---Hello!
}

func main() {
	var v int = 50
	fmt.Printf("%c", v)
}
