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
	s = strings.Trim(s, "\n ")
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
