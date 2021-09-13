package gonet

import (
	"fmt"
	"net"
	"os"
	"testing"
	"time"
)

func TestLookupIP(t *testing.T) {
	hosts := []string{
		"www.baidu.com",
		"www.cn.bing.com",
		"www.google.com",
	}

	for _, host := range hosts {
		start := time.Now()
		ips, err := net.LookupIP(host)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		fmt.Printf("Host: %s; Time usage:%v, ips:%v\n", host, time.Since(start), ips)
	}
}

func TestLookupAddr(t *testing.T) {
	addrs := []string{
		"114.114.114.114",
		"202.96.134.133",
		"223.5.5.5",
	}

	for _, addr := range addrs {
		start := time.Now()
		names, err := net.LookupAddr(addr)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		fmt.Printf("Addr: %s; Time usage:%v, names:%v\n", addr, time.Since(start), names)
	}
}

func TestResolveAddr(t *testing.T) {
	host := "www.baidu.com"
	addr, err := net.ResolveIPAddr("", host)
	//net.ResolveTCPAddr()
	//net.ResolveUDPAddr()
	//net.ResolveUnixAddr()

	if err != nil {
		fmt.Println("Error:", err.Error())
		return
	}
	fmt.Fprintf(os.Stdout, "%s IP: %s Network: %s Zone: %s\n", addr.String(), addr.IP, addr.Network(), addr.Zone)

}
