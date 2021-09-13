package text

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/bitly/go-simplejson"
)

type Server struct {
	ServerName string `json:"serverName"`
	ServerIp   string `json:"serverIp"`
}

type ServerSlice struct {
	Servers []Server `json:"servers"`
}

func getJSONInfo() {
	file, err := os.OpenFile("./server.json", 0, os.FileMode(os.O_CREATE))
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	var s ServerSlice
	json.Unmarshal(data, &s)
	fmt.Println(s)
}

func getJSONInfoUsingInterface() {
	b := []byte(`{"Name":"Wednesday","Age":6,"Parents":["Gomez","Morticia"]}`)

	var f interface{}
	err := json.Unmarshal(b, &f)
	if err != nil {
		fmt.Printf("err: %v", err)
	}
	fmt.Println(f)

	content := f.(map[string]interface{})
	for key, value := range content {
		switch vv := value.(type) {
		case string:
			fmt.Println(key, "is string", vv)
		case int:
			fmt.Println(key, "is int", vv)
		case float64:
			fmt.Println(key, "is float64", vv)
		case []interface{}:
			fmt.Println(key, "is an array:")
			for i, u := range vv {
				fmt.Println(i, u)
			}
		default:
			fmt.Println(key, "is of a type I don't know how to handle")
		}
	}

	b = []byte(`{
		"test": {
			"array": [1, "2", 3],
			"int": 10,
			"float": 5.150,
			"bignum": 9223372036854775807,
			"string": "simplejson",
			"bool": true
		}
	}`)
	js, err := simplejson.NewJson(b)
	testObj := js.Get("test")
	arr, _ := testObj.Get("array").Array()
	fmt.Println(arr)
	i, _ := testObj.Get("int").Int()
	f64, _ := testObj.Get("bool").Bool()
	fmt.Println(i, f64)
}

func setJSONInfo() {
	type Server struct {
		ServerName string `json:"serverName"`
		ServerIP   string `json:"serverIP"`
	}

	type ServerSlice struct {
		Servers []Server
	}

	var s ServerSlice
	s.Servers = append(s.Servers, Server{ServerName: "Shanghai_VPN", ServerIP: "127.0.0.1"})
	s.Servers = append(s.Servers, Server{ServerName: "Beijing_VPN", ServerIP: "127.0.0.2"})
	b, err := json.Marshal(s)
	if err != nil {
		fmt.Println("json err:", err)
	}
	fmt.Println(string(b))
}
