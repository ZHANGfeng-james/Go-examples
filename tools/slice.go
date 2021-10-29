package tools

import (
	"fmt"
	"log"
	"reflect"
)

// SliceInfo log slice infomation of len/cap/element ...
func SliceInfo(variable string, v interface{}) string {
	// 入参是 interface{} --> slice, for example []int
	value := reflect.ValueOf(v)
	if value.Kind() != reflect.Slice {
		log.Fatalf("sliceInfo error! %s interface value is not reflect.Slice", variable)
	}
	// ele := value.Slice(0, value.Len())
	return fmt.Sprintf("%8s slice len():%d, cap():%d, ele:%v", variable, value.Len(), value.Cap(), value)
}
