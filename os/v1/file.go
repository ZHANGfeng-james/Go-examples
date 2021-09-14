package main

import (
	"io"
	"log"
	"os"
	"strconv"
)

func main() {
	// Test filename
	file, err := os.Open("./20220914114008_M008.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// copy file
	num := 20220914114008
	suffix := "_M008.jpg"
	endNum := num + 10
	for num += 1; num < endNum; num++ {
		name := strconv.Itoa(num) + suffix
		// create a new file
		dstFile, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			log.Fatal(err)
			return
		}
		defer dstFile.Close()

		io.Copy(dstFile, file)
		// must be call seek to set the offset for next read!
		file.Seek(0, io.SeekStart)
	}
}
