package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

type FileInfo struct {
	name    string
	size    int64
	modTime time.Time
}

func (info FileInfo) String() string {
	return fmt.Sprintf(`
	name: %s
	size:%d(Byte)
	modTime:%v
	`, info.name, info.size, info.modTime)
}

const (
	N = 3
)

func copyFileUseioCopy() {
	// Test filename
	file, err := os.Open("./copy_file_test.go")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}

	fileInfo := FileInfo{
		name:    info.Name(),
		size:    info.Size(),
		modTime: info.ModTime(),
	}
	fmt.Println(fileInfo)

	// copy file
	num := 20220914114008
	suffix := "_M008.txt"
	endNum := num + N
	for num += 1; num < endNum; num++ {
		name := strconv.Itoa(num) + suffix
		// create a new file
		dstFile, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			log.Fatal(err)
			return
		}
		defer func(index int) {
			fmt.Println("file close, ", index)
			// run after all sentences before copyFileUseioCopy() over
			dstFile.Close()
		}(num)

		io.Copy(dstFile, file)
		// must be call seek to set the offset for next read!
		file.Seek(0, io.SeekStart)
		fmt.Println("for loop, ", num)
	}
	fmt.Println("over")
}
