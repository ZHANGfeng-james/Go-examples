package v4

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func TestGetter(t *testing.T) {
	var f Getter = GetterFunc(func(key string) ([]byte, error) { // 函数类型的初始化，并赋值给一个接口
		return []byte(key), nil
	})
	expect := []byte("key")
	if v, _ := f.Get("key"); !reflect.DeepEqual(v, expect) {
		t.Fatal("callback failed")
	}
}

func TestGroup(t *testing.T) {
	NewGroup("score", 2<<10, GetterFunc(func(key string) (bytes []byte, err error) {
		return
	}))

	if group := GetGroup("score"); group == nil || group.name != "score" {
		t.Fatal("create new group failed")
	}
	if group := GetGroup("score" + "xxx"); group != nil {
		t.Fatal("get a error group")
	}
}

func TestGet(t *testing.T) {
	loadCount := make(map[string]int)
	gee := NewGroup("school score", 2<<10, GetterFunc(func(key string) ([]byte, error) {
		log.Println("[SlowDB] search key:", key)
		if v, ok := db[key]; ok {
			if _, ok := loadCount[key]; !ok {
				loadCount[key] = 0
			}
			// 表示从db中加载数据的次数
			loadCount[key]++
			return []byte(v), nil
		}
		return nil, fmt.Errorf("%s not exist", key)
	}))
	for k, v := range db {
		if view, err := gee.Get(k); err != nil || view.String() != v {
			t.Log(err, view.String())
			t.Fatalf("failed to get value of %s", k)
		}
		if _, err := gee.Get(k); err != nil || loadCount[k] > 1 {
			t.Fatalf("cache %s miss", k)
		}
	}
	if view, err := gee.Get("unknown"); err == nil {
		t.Fatalf("the value of unknow should be empty, but %s got", view)
	}
}
