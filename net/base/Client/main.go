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
