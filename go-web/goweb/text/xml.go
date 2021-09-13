package text

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)

type RecurlyServer struct {
	XMLName xml.Name `xml:"servers"`
	Version string   `xml:"version,attr"`
	Svs     []server `xml:"server"`
}

type server struct {
	XMLName    xml.Name `xml:"server"`
	ServerName string   `xml:"serverName"`
	ServerIp   string   `xml:"serverIP"`
}

func getXMLInfo() {
	file, err := os.OpenFile("./server.xml", 0, os.FileMode(os.O_CREATE))
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

	v := RecurlyServer{}
	err = xml.Unmarshal(data, &v)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	fmt.Println(v)
}

func setXMLInfo() {
	type server struct {
		ServerName string `xml:"serverName"`
		ServerIp   string `xml:"serverIP"`
	}
	type RecurlyServer struct {
		XMLName xml.Name `xml:"servers"`
		Version string   `xml:"version,attr"`
		Svs     []server `xml:"server"`
	}

	v := &RecurlyServer{
		Version: "1",
	}

	v.Svs = append(v.Svs, server{"Shanghai_VPN", "127.0.0.1"})
	v.Svs = append(v.Svs, server{"Beijing_VPN", "127.0.0.2"})

	output, err := xml.MarshalIndent(v, "  ", "	")
	if err != nil {
		fmt.Printf("err: %v.\n", err)
	}

	os.Stdout.Write([]byte(xml.Header))
	os.Stdout.Write(output)
}
