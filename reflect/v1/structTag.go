package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Data struct {
	Labels     []string `http:"labels"`
	MaxResults int      `http:"max"`
	Exact      bool     `http:"exact"`
}

func test() {
	data := Data{
		MaxResults: 10,
	}

	// 模拟 req.ParseForm()，解析出 URL 中的 key-value
	type Values map[string][]string
	form := make(Values)
	form["labels"] = []string{"a", "b", "c"}
	form["exact"] = []string{"truee"}
	// URL 中的 key-value 对应解析到 Data 结构体变量中
	Unpack(form, &data)

	fmt.Println(data)
}

func Unpack(form map[string][]string, ptr interface{}) error {
	valueOf := reflect.ValueOf(ptr) // the pointer variable
	ele := valueOf.Elem()           // the struct variable
	fmt.Println(valueOf, fmt.Sprintf("%T", valueOf))
	fmt.Println(ele, fmt.Sprintf("%T", ele))

	// 获取 Data 结构体的 Field 信息
	fields := make(map[string]reflect.Value)
	for i := 0; i < ele.NumField(); i++ {
		// 获取 ele 的类型值 reflect.Type
		fieldInfo := ele.Type().Field(i)
		tag := fieldInfo.Tag
		name := tag.Get("http")
		if name == "" {
			name = strings.ToLower(fieldInfo.Name)
		}
		fields[name] = ele.Field(i) // labels --> <[]string Value>
	}
	// map[exact:<bool Value> labels:<[]string Value> max:<int Value>]
	fmt.Println(fields)

	for name, values := range form {
		// 取到 name 对应的 reflect.Value
		f := fields[name]
		if !f.IsValid() {
			continue
		}

		for _, value := range values {
			if f.Kind() == reflect.Slice {
				// reflect.Slice 的元素类型，依据该类型创建一个ptr
				elem := reflect.New(f.Type().Elem()).Elem()
				if err := populate(elem, value); err != nil {
					return fmt.Errorf("%s: %v", name, err)
				}
				f.Set(reflect.Append(f, elem))
			} else {
				if err := populate(f, value); err != nil {
					return fmt.Errorf("%s: %v", name, err)
				}
			}
		}
	}

	return nil
}

func populate(v reflect.Value, value string) error {
	switch v.Kind() {
	case reflect.String:
		// reflect.Value 其底层值设置为 value
		v.SetString(value)
	case reflect.Int:
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		v.SetInt(i)
	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		v.SetBool(b)
	default:
		return fmt.Errorf("unsupported kind %s", v.Type())
	}
	return nil
}
