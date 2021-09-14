package goioutil

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestReadFile(t *testing.T) {
	binaryFile := "./ubuntu-16.04.6-desktop-i386.iso"
	textFile := "./README.text"
	fmt.Println(binaryFile, textFile)

	file, err := os.Open(textFile)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("Open File success!", file.Name())
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		log.Fatal(err.Error())
	}

	size := float32(info.Size()) / 1024 / 1024 / 1024
	fmt.Println("Get file info success!", size, ("GB"), "<--", info.Size())

	var buf []byte
	buf = make([]byte, 200)

	for {
		n, err := file.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err.Error())
		}
		fmt.Println(n)
	}
}

func ReadAll(file *os.File) {
	// 获取 File 实例，且其实现了 io.Reader 接口
	fmt.Println("Start to read File!")
	contentBytes, err := ioutil.ReadFile(file.Name())
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("Read End!", len(contentBytes))
}
