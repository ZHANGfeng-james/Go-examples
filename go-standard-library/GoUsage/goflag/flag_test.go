package goflag

import (
	"flag"
	"fmt"
	"os"
	"testing"
)

func TestFlagNormal(t *testing.T) {
	if !flag.Parsed() {
		flag.Parse()
	}
	flags := flag.Args()
	fmt.Println(flags)

	all := os.Args
	fmt.Println(all, "size:", len(os.Args))

	args := os.Args[len(os.Args)-1]
	fmt.Println(args)

	var nFlag = flag.Int("name", 1234, "help message for flag name")
	fmt.Println(*nFlag)
}
