package usage

import (
	"flag"
	"fmt"
	"log"
)

func ParseFlagSet() {
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
