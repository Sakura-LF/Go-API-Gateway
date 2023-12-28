package http

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var (
	//代理服务器地址
	proxyAddrDemo = "127.0.0.1:8081"
)

// ReverseProxySever 使用 ReverProxy 实现反向代理
func ReverseProxySeverDemo() {
	// 解析下游服务器的真实地址
	realServer := "http://127.0.0.1:8001"
	serverUrl, err := url.Parse(realServer)
	if err != nil {
		log.Fatalln("parse fail,")
	}
	proxy := httputil.NewSingleHostReverseProxy(serverUrl)
	log.Println("Starting proxy Server:", proxyAddrDemo)
	// proxy 就相当于 handle
	http.ListenAndServe(proxyAddrDemo, proxy)
}
