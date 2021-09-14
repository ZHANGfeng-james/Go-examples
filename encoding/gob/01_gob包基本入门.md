encoding/gob 包提供了一种**序列化**和**反序列化** Go **数据类型**的方法。

encoding/gob 包中有 2 个基本实体：

* Encoder：编码器，用于管理**类型和数据信息**传输到连接的另一端，是并发安全的；
* Decoder：解码器，用于管理从连接的远端接收来的**类型和数据信息**，是并发安全的。



~~~go
package gob

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"testing"
)

type P struct {
	X, Y, Z int
	Name    string
}

type Q struct {
	X, Y *int32
	Name string
}

func TestEndingGob(t *testing.T) {
	var buf bytes.Buffer

	encoder := gob.NewEncoder(&buf)
	decoder := gob.NewDecoder(&buf)

	err := encoder.Encode(P{
		3, 4, 5,
		"Michoi",
	})
	if err != nil {
		log.Fatal("encode error:", err)
	}

	var instance Q
	err = decoder.Decode(&instance)
	if err != nil {
		log.Fatal("decode error:", err)
	}
	fmt.Printf("%q: {%d, %d}\n", instance.Name, *instance.X, *instance.Y)
}
"Michoi": {3, 4}
~~~

