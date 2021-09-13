package dir

import (
	"fmt"
	"os"
)

func document() {
	os.Mkdir("astaxie", 0777)
	os.MkdirAll("astaxie/test1/teste2", 0777)
	err := os.Remove("astaxie")
	if err != nil {
		fmt.Println(err)
	}

	os.RemoveAll("astaxie")
}
