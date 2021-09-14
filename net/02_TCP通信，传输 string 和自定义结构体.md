本文通过 TCP 通信在 Server 和 Client 端传输 string 和自定义结构体：

~~~bash
                                客户端                                        服务端

                                待发送结构体                                解码后结构体
                                testStruct结构体                            testStruct结构体
                                    |                                             ^
                                    V                                             |
                                gob编码       ---------------------------->     gob解码
                                    |                                             ^
                                    V                                             |   
                                   发送     ============网络=================    接收
~~~

其中自定义结构体实例通过 ending/gob 包编解码实现二进制流在网络中的传输。

Client 端代码：

~~~go
package main

import (
	"bufio"
	"encoding/gob"
	"log"
	"net"
	"strconv"

	"github.com/pkg/errors"
)

const (
	// Port is the port number that the server listens to.
	Port = ":61000"
)

type complexData struct {
	N int
	S string
	M map[string]int
	P []byte
	C *complexData
}

func Open(addr string) (*bufio.ReadWriter, error) {
	log.Println("Dial " + addr)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, errors.Wrap(err, "Dialing "+addr+" failed")
	}
	return bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn)), nil
}

func main() {
	testStruct := complexData{
		N: 23,
		S: "string data",
		M: map[string]int{"one": 1, "two": 2, "three": 3},
		P: []byte("abc"),
		C: &complexData{
			N: 256,
			S: "Recursive structs? Piece of cake!",
			M: map[string]int{"01": 1, "10": 2, "11": 3},
		},
	}

	ip := "127.0.0.1"

	rw, err := Open(ip + Port)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Client: Failed to open connection to "+ip+Port).Error())
	}
	// Send a STRING request.
	// Send the request name.
	// Send the data.
	log.Println("Send the string request.")
	cmd := "STRING\n"
	n, err := rw.WriteString(cmd)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Could not send the STRING request ("+strconv.Itoa(n)+" bytes written)").Error())
	}
	n, err = rw.WriteString("Additional data.\n")
	if err != nil {
		log.Fatal(errors.Wrap(err, "Could not send additional STRING data ("+strconv.Itoa(n)+" bytes written)").Error())
	}
	log.Println("Flush the buffer.")
	err = rw.Flush()
	if err != nil {
		log.Fatal(errors.Wrap(err, "Flush failed.").Error())
	}

	// Read the reply.
	log.Println("Read the reply.")
	response, err := rw.ReadString('\n')
	if err != nil {
		log.Fatal(errors.Wrap(err, "Client: Failed to read the reply: '"+response+"'").Error())
	}
	log.Println("STRING request: got a response:", response)

	log.Println("Send a struct as GOB")
	log.Printf("Outer complexData struct: \n%#v\n", testStruct)
	log.Printf("Inner complexData struct: \n%#v\n", testStruct.C)

	n, err = rw.WriteString("GOB\n")
	if err != nil {
		log.Fatal(errors.Wrap(err, "Could not write GOB data ("+strconv.Itoa(n)+" bytes written)").Error())
	}
	enc := gob.NewEncoder(rw)
	err = enc.Encode(testStruct)
	if err != nil {
		log.Fatal(errors.Wrapf(err, "Encode failed for struct: %#v", testStruct).Error())
	}
	err = rw.Flush()
	if err != nil {
		log.Fatal(errors.Wrap(err, "Flush failed.").Error())
	}

	cmd = "STRING\n"
	n, err = rw.WriteString(cmd)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Could not send the STRING request ("+strconv.Itoa(n)+" bytes written)").Error())
	}
	n, err = rw.WriteString("Goodbye!\n")
	if err != nil {
		log.Fatal(errors.Wrap(err, "Could not send additional STRING data ("+strconv.Itoa(n)+" bytes written)").Error())
	}
	log.Println("Flush the buffer.")
	err = rw.Flush()
	if err != nil {
		log.Fatal(errors.Wrap(err, "Flush failed.").Error())
	}
	// Read the reply.
	log.Println("Read the reply.")
	response, err = rw.ReadString('\n')
	if err != nil {
		log.Fatal(errors.Wrap(err, "Client: Failed to read the reply: '"+response+"'").Error())
	}
	log.Println("STRING request: got a response:", response)

	log.Println("Client done.")
}
~~~

Server 端代码：

~~~go
package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

const (
	Port = ":61000"
)

// A struct with a mix of fields, used for the GOB example.
type complexData struct {
	N int
	S string
	M map[string]int
	P []byte
	C *complexData
}

type HandleFunc func(*bufio.ReadWriter)

type Endpoint struct {
	listener net.Listener
	handler  map[string]HandleFunc

	m sync.RWMutex
}

func NewEndPoint() *Endpoint {
	return &Endpoint{
		handler: map[string]HandleFunc{},
	}
}

func (e *Endpoint) AddHandleFunc(name string, f HandleFunc) {
	e.m.Lock()
	defer e.m.Unlock()
	e.handler[name] = f
}

func (e *Endpoint) Listen() error {
	var err error
	e.listener, err = net.Listen("tcp", Port)
	if err != nil {
		return errors.Wrapf(err, "Unable to listen on port %s\n", Port)
	}
	log.Println("Listen on", e.listener.Addr().String())
	for {
		log.Println("Waiting for Accept a connection request...")
		conn, err := e.listener.Accept()
		if err != nil {
			log.Println("Failed accepting a connection request:", err)
			continue
		}
		log.Println("Handle incoming message.")
		e.handleMessage(conn)
	}
}

func (e *Endpoint) handleMessage(conn net.Conn) {
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	defer conn.Close()

	for {
		cmd, err := rw.ReadString('\n')
		switch {
		case err == io.EOF:
			log.Println("Reached EOF - close this connection.\n   ---")
			return
		case err != nil:
			log.Println("\nError reading command. Got: '"+cmd+"'\n", err)
			return
		}

		fmt.Printf("%q\n", []byte(cmd))
		cmd = strings.Trim(cmd, "\n")

		e.m.Lock()
		handler, ok := e.handler[cmd]
		e.m.Unlock()
		if !ok {
			log.Println("Command '" + cmd + "' is not registered.")
			return
		}
		handler(rw)
	}

}

func handleString(rw *bufio.ReadWriter) {
	// Receive a string.
	log.Print("Receive STRING message:")
	s, err := rw.ReadString('\n')
	if err != nil {
		log.Println("Cannot read from connection.\n", err)
	}
	s = strings.Trim(s, "\n")
	log.Println(s)
	_, err = rw.WriteString("Thank you.\n")
	if err != nil {
		log.Println("Cannot write to connection.\n", err)
	}
	err = rw.Flush()
	if err != nil {
		log.Println("Flush failed.", err)
	}
}

func handleGob(rw *bufio.ReadWriter) {
	log.Print("Receive GOB data:")
	var data complexData
	// Create a decoder that decodes directly into a struct variable.
	dec := gob.NewDecoder(rw)
	err := dec.Decode(&data)
	if err != nil {
		log.Println("Error decoding GOB data:", err)
		return
	}
	// Print the complexData struct and the nested one, too, to prove
	// that both travelled across the wire.
	log.Printf("Outer complexData struct: \n%#v\n", data)
	log.Printf("Inner complexData struct: \n%#v\n", data.C)
}

func main() {
	endpoint := NewEndPoint()

	endpoint.AddHandleFunc("STRING", handleString)
	endpoint.AddHandleFunc("GOB", handleGob)

	err := endpoint.Listen()
	if err != nil {
		log.Println("Error:", errors.WithStack(err))
	}
}
~~~

