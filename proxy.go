package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
)

func handleClientRequest(client net.Conn) {
	if client == nil {
		return
	}
	defer client.Close()

	var b [1024]byte
	n, err := client.Read(b[:])
	if err != nil {
		log.Panic(err.Error())
		return
	}
	var method, host, address string
	fmt.Sscanf(string(b[:bytes.IndexByte(b[:], '\n')]), "%s%s", &method, &host)
	hostPortURL, err := url.Parse(host)
	if err != nil {
		log.Panicln("err_0", err)
		return
	}
	log.Println("Scheme.port", hostPortURL.Port())
	address = hostPortURL.Host + ":80"
	server, err := net.Dial("tcp", address)

	if err != nil {
		log.Println(err)
		return
	}
	if method == "CONNECT" {
		fmt.Fprint(client, "HTTP/1.1 200 Connection established\r\n")
	} else {
		server.Write(b[:n])
	} //进行转发
	go io.Copy(server, client)
	io.Copy(client, server)
	//	var method, host, address string
}

// 启动
func main() {
	// 设置默认启动监控的ip和端口
	var addr = flag.String("addr", ":8097", "The addr of the application.")
	flag.Parse() // 解析参数
	fmt.Println(*addr)
	l, err0 := net.Listen("tcp", *addr)
	if err0 != nil {
		log.Panic(err0)
	}
	for {
		client, err1 := l.Accept()
		if err1 != nil {
			log.Panic(err1)
		}
		go handleClientRequest(client)
	}
}
