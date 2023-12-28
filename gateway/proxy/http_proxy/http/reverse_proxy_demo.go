package http

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

func CustomerReverseProxySever() {
	var port = "8081" // 当前代理服务器端口
	http.HandleFunc("/", handler)
	fmt.Println("反向代理服务器启动: " + port)
	http.ListenAndServe(":"+port, nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("经过代理服务器地址")
	// 1.解析下游服务器地址，更改请求地址
	// 被代理的下游真实服务器地址.应该通过一定的负载均衡算法获取
	realServer, err := url.Parse(proxyAddr)
	r.URL.Scheme = realServer.Scheme // http
	r.URL.Host = realServer.Host     // 127.0.0.1:8001

	// 2.请求下游(真实服务器)，并获取返回内容
	transport := http.DefaultTransport
	resp, err := transport.RoundTrip(r) // 得到下游服务器响应
	defer resp.Body.Close()
	if err != nil {
		log.Print(err)
		return

	}
	// 3.把下游请求内容做一些处理，然后返回给上游(客户端)
	for k, vv := range resp.Header { // 修改上游响应头
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}
	bufio.NewReader(resp.Body).WriteTo(w) // 将下游响应体写回上游客户端
}
