package main

import (
	"flag"
	"log"
	"net/http"
)

// 启动
func main() {
	// 设置默认启动监控的ip和端口
	var addr = flag.String("addr", ":8097", "The addr of the application.")
	flag.Parse() // 解析参数

	log.Println("Starting proxy server on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
