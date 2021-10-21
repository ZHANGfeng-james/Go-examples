package usage

import (
	"flag"
	"log"
	"os"
)

func ParseFlag() {
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

func ReadFromVariable() {
	var nFlag = flag.Int("i", 10, "flag param n value")
	*nFlag = 20
	// 直接就可以修改 Flag 的值
	flag := flag.CommandLine.Lookup("i")
	log.Printf("read:%s", flag.Value.String())

	// 一般情况下是，CLI 作为数据源（输入），经过 flag 包的解析，在程序中可获得 CLI 的输入
	// 上面这种情况，并不实用
}
