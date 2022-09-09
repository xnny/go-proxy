package main

// ex $ curl --proxy http://127.0.0.1:12345 -L https://www.google.com

import (
	"bytes"
	"io"
	"log"
	"net"
	"net/url"
	"strconv"
	"strings"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	l, err := net.Listen("tcp", ":8097")
	if err != nil {
		log.Panic(err)
	}
	defer func(l net.Listener) {
		err := l.Close()
		if err != nil {
			log.Panic(err)
		}
	}(l)

	for {
		client, err := l.Accept()
		if err != nil {
			log.Panic(err)
		}
		go handleClientRequest(client)
	}
}

func handleClientRequest(client net.Conn) {
	if client == nil {
		return
	}
	defer func(client net.Conn) {
		err := client.Close()
		if err != nil {
			log.Panic(err)
		}
	}(client)

	var b [1024]byte
	n, err := client.Read(b[:])
	if err != nil {
		log.Println(err)
		return
	}

	var method, host, address string

	ni := bytes.IndexByte(b[:], '\n')
	if ni != -1 {
		// 避免url 过长崩溃
		sl := strings.Split(string(b[:ni]), " ")
		if len(sl) >= 2 {
			method, host = sl[0], sl[1]
			if _, err := strconv.Atoi(string(host[0])); err == nil {
				// fixed  https://github.com/golang/go/issues/19297
				host = "//" + host
			}
		}
	}

	hostPortURL, err := url.Parse(host)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(hostPortURL.Scheme + "|" + hostPortURL.Opaque)

	if len(hostPortURL.Opaque) > 0 { // 如果是带证书请求
		address = hostPortURL.Scheme + ":" + hostPortURL.Opaque
	} else {
		if strings.Index(hostPortURL.Host, ":") == -1 { // 简单的 http 请求
			address = hostPortURL.Host + ":80"
		} else {
			address = hostPortURL.Host
		}
	}

	// 获得了请求的host和port，就开始拨号吧
	server, err := net.Dial("tcp", address)
	if err != nil {
		log.Println(err)
		return
	}
	if method == "CONNECT" {
		_, err := client.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))
		if err != nil {
			log.Println(err)
		}
	} else {
		_, err := server.Write(b[:n])
		if err != nil {
			log.Println(err)
		}
	}

	//进行转发
	go func() {
		_, err := io.Copy(server, client)
		if err != nil {
			log.Println(err)
		}
	}()
	_, err = io.Copy(client, server)
	if err != nil {
		log.Println(err)
	}
}
