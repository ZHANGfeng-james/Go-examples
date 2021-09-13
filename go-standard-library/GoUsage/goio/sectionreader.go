package main

import (
	"io"
	"log"
	"os"
)

func main() {
	file, _ := os.Open("./A20_Apex.apk")
	defer file.Close()
	fileinfo, err := file.Stat()
	if err != nil {
		log.Fatal(err.Error())
	}

	r := io.NewSectionReader(file, 0, fileinfo.Size())
	reader := sectionReadCloser{r}
	defer reader.Close()

	out, err := os.Create("./copy.apk")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer out.Close()
	// 将 src 拷贝到 out 中
	io.Copy(out, reader)
}

type sectionReadCloser struct {
	*io.SectionReader
}

func (rc sectionReadCloser) Close() error {
	return nil
}
