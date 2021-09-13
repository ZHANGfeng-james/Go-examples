package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	TimeParseFormat()
}

func ParseDuration() {
	duration, err := time.ParseDuration("5m")
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(duration)
}

func TimeParseFormat() {
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(location)

	inputTime := "2029-09-04 12:02:33"
	layout := "2006-01-02 15:04:05"

	//t, _ := time.Parse(layout, inputTime)
	t, _ := time.ParseInLocation(layout, inputTime, location)

	dateTime := time.Unix(t.Unix(), 0).In(location).Format(layout)
	fmt.Printf("输入时间：%s, 输出时间:%s\n", inputTime, dateTime)
}
