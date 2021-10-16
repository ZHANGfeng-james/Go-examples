package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	parseFlagSet()
}

func parseFlagSet() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// 实际的作用就是取到 flag.Args，使用默认的 FlagSet
	flag.Parse()

	var name string

	goCmd := flag.NewFlagSet("go", flag.ExitOnError)       // 创建 name 为 go 的 FlagSet
	goCmd.StringVar(&name, "name", "Go语言", "help message") // goCmd 这个 FlagSet 中为 name 变量预解析参数标识符

	javaCmd := flag.NewFlagSet("java", flag.ExitOnError)
	javaCmd.StringVar(&name, "name", "Java语言", "help message") // javaCmd 这个 FlagSet 中为 name 变量预解析参数标识符

	// 取到 os.Args[1:] 的命令行参数
	args := flag.Args()
	if len(args) <= 0 {
		return
	}
	fmt.Printf("%d, %v\n", len(args), args)

	// 匹配到对应的子命令
	switch args[0] {
	case "go":
		// 解析接下来的命令行参数
		_ = goCmd.Parse(args[1:])
	case "java":
		_ = javaCmd.Parse(args[1:])
	}

	fmt.Printf("name=%q\n", name)
}

func parseFlag() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var nFlag = flag.Int("i", 10, "flag param n value")

	var nameFlag string
	flag.StringVar(&nameFlag, "s", "Katyusha", "name of company")

	var boolFlag bool
	flag.BoolVar(&boolFlag, "b", false, "whether or not")

	if !flag.Parsed() {
		flag.Parse()
	}
	log.Printf("nFlag:%d, nameFlag:%s, boolFlag:%v", *nFlag, nameFlag, boolFlag)

	args := os.Args
	log.Println(args)

	flagArgs := flag.Args()
	log.Println(flagArgs)
}

func parseArgv() {
	flags := flag.Args()
	fmt.Println(flags)

	all := os.Args
	fmt.Println(all, "size:", len(os.Args))

	args := os.Args[len(os.Args)-1]
	fmt.Println(args)

	var nFlag = flag.Int("name", 1234, "help message for flag name")
	fmt.Println(*nFlag)
}
