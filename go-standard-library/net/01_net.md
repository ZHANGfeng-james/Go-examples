

# 1 net 标准库

Go 中的 net 包，提供了面向网络 I/O 的**可移植性接口**，包括：TCP/IP、UDP、域名解析，以及 Unix 域套接字（Socket）。虽然 net 包提供了**网络底层原语**的访问方式，但大多数客户端仅需要包括 Dial、Listen 和 Accept 函数，以及和 Conn、Listener 相关的基本接口。

crypto/tls 包使用了和 net 相同的接口，以及相同的 Dial 和 Listen 函数。

![](./img/Snipaste_2021-07-08_09-03-40.png)

此处特别指出：**网络底层原语**—— low-level networking primitives，就要提到**计算机网络通信模型** OSI/RM（Open System Interconnection Reference Model），由全世界的计算机企业以此为标准开发和建立计算机网络，最终扩展到各种不同类型的终端设备的组网问题。由此实现互联网建立之初的目的：**把全世界的计算机连接起来，实现信息共享**！

**整个模型**是这样的：

| 层级 |    名称    |           主要功能           |                    主要设备和协议                    |
| :--: | :--------: | :--------------------------: | :--------------------------------------------------: |
|  7   |   应用层   |      实现具体的应用功能      | POP3、FTP、HTTP、Telnet、SMTP、DHCP、TFTP、SNMP、DNS |
|  6   |   表示层   | 数据的格式与表达、加密、压缩 |                                                      |
|  5   |   会话层   |     建立、管理和终止会话     |                                                      |
|  4   |   传输层   |          端到端连接          |                       TCP、UDP                       |
|  3   |   网络层   |    **分组**传递和路由选择    |   三层交换机、路由器（ARP、RARP、IP、ICMP、IGMP）    |
|  2   | 数据链路层 |   传输以**帧**为单位的信息   |     网桥、交换机、网卡（PPTP、L2TP、SLIP、PPP）      |
|  1   |   物理层   |          二进制传输          |    中继器和集线器（ISO2110，IEEE802.1，EEE802.2）    |

它的工作规则是：每层向上层提供服务，同时使用下层提供的服务，且不可跨层。

对于网络编程来说，需要理解：

* **传输层**：为上层协议提供端到端的可靠和透明的数据传输服务，包括处理**差错控制**和**流量控制**等问题。该层向高层屏蔽了下层数据通信的细节，使高层用户看到的只是在**两个传输实体间的一条主机到主机的、可由用户控制和设定的、可靠的数据通路**；
* **会话层**：负责**建立、管理和终止**表示层实体之间的**通信会话**，该层的通信由不同设备中的应用程序之间的服务**请求**和**响应**组成；
* 表示层：提供各种用于应用层**数据的编码和转换**功能，确保一个系统的应用层发送的数据能被另一个系统的应用层识别；
* 应用层：为计算机用户提供应用接口，也为用户直接提供各种网络服务。

和模型直接相关的是各层使用的**协议**，而协议是什么？是**网络上各种计算设备之间进行交流的语言**。为了实现设备间的通信，就要确保通信双方实体完成通信所必需遵循的规则和约定，确保数据单元使用的格式，信息单元信息与含义，信息发送和接收的时序满足确定且一致的标准规范，这个标准规范就是**协议**。而协议的实现考虑的就是：**语义、语法和时序**，这个过程可类比到我们日常的口头、书信沟通活动。

有了计算机网络通信模型的基础，我们就可以**正式进入网络编程的世界**了！

# 2 总览

通过 Dial 函数连接 Server，相当于是**创建 Client 应用程序**：

~~~go
conn, err := net.Dial("tcp", "golang.org:80")
if err != nil {
	// handle error
}
fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
status, err := bufio.NewReader(conn).ReadString('\n')
// ...
~~~

使用 Listen **创建 Server 应用程序**：

~~~go
ln, err := net.Listen("tcp", ":8080")
if err != nil {
	// handle error
}
for {
	conn, err := ln.Accept()
	if err != nil {
		// handle error
	}
	go handleConnection(conn)
}
~~~

从 net 的整个内容来看，是和 net/http 包有很大区别的：net 包提供的是网络通信的底层原语，比如：TCP/IP、UDP、域名解析，以及 Unix 域套接字（Socket）。而 net/http 包则提供的是和这个 HTTP 应用层协议相关的封装结构，比如：Request、Response、Header、Cookie 等。

# 3 域名解析

关于主机域名解析，有间接的解析比如 Dial，以及直接的解析 LookupHost 和 LookupAddr，可根据操作系统的不同，使用不同的函数。

**Unix 系统中**，有 2 种可选方式用于解析域名。可以使用纯 Go 解析器，发送 DNS 请求给 /etc/resolv.conf 中标记的服务器；或者，调用基于 cgo 的解析器，会调用 C 库的函数，比如 getaddrinfo 和 getnameinfo；

默认情况下，会使用纯 Go 解析器，因为一个阻塞的 DNS 查询会消耗一个 goroutine。但是一个阻塞的 C 调用会消耗一个操作系统线程。当 cgo 可用时，基于 cgo 的解析器会在如下情况下被使用：

1. OS X 中，操作系统不允许发送 DNS 请求；
2. LOCALDOMAIN 环境变量被设置；
3. RES_OPTIONS 或 HOSTALIASES 环境变量不为空；
4. ASR_CONFIG 环境变量不为空；
5. /etc/resolv.conf 或 /etc/nsswitch.conf 需要使用的特性 Go 解析器没有实现；
6. 待查找的域名在 .local 中不存在或者是一个 mDNS 的名称。

解析器的使用选择可以被如下环境变量覆写：

~~~bash
export GODEBUG=netdns=go    # force pure Go resolver
export GODEBUG=netdns=cgo   # force cgo resolver
~~~

在编译 Go 源代码时，也可以通过设置 netgo 或 netcgo 的编译标签强制设置解析器的行为。

在 Plan 9 操作系统中，解析器会返回 /net/cs 和 /net/dns；在 Windows 系统中，解析器会使用 C 库中的函数，比如 GetAddrInfo 和 DnsQuery。













`func LookupIP(host string) ([]IP, error)`：依据主机地址解析到 IP 地址（**域名解析服务 DNS**）。

~~~go
package gonet

import (
	"fmt"
	"net"
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
Host: www.baidu.com; Time usage:15.9908ms, ips:[14.215.177.38 14.215.177.39]
Host: www.cn.bing.com; Time usage:12.9937ms, ips:[204.79.197.200]
Host: www.google.com; Time usage:0s, ips:[103.200.30.143]
~~~

从实际来看，DNS 域名解析服务是一个耗时的操作，一般情况下是 ms 级别的操作。

`func LookupAddr(addr string) (names []string, err error)`：对给定地址指向反向查找，找到对应的 DNS 地址列表

~~~go
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
Addr: 114.114.114.114; Time usage:11.9934ms, names:[public1.114dns.com.]
Addr: 202.96.134.133; Time usage:5.999ms, names:[ns.szptt.net.cn.]
Addr: 223.5.5.5; Time usage:20.9834ms, names:[public1.alidns.com.]
~~~

有下列关于**网络地址的解析函数**：

1. `   func ResolveIPAddr(network, address string) (*IPAddr, error)`：

2. `func ResolveTCPAddr(network, address string) (*TCPAddr, error)`：network 必须是 TCP

3. `func ResolveUDPAddr(network, address string) (*UDPAddr, error)`：

4. `func ResolveUnixAddr(network, address string) (*UnixAddr, error)`：

