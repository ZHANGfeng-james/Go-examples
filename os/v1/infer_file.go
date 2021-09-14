package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func inferRootDir() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Println("cwd:", cwd)

	var infer func(d string) string
	infer = func(d string) string {
		if exists(d + "/configs") {
			return d
		}
		return infer(filepath.Dir(d))
	}

	RootDir := infer(cwd)
	fmt.Println(RootDir)
}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
