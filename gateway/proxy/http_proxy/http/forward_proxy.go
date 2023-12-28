package http

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
)

func ForwardProxyServer1() {
	http.Handle("/", &Proxy{})
	// 代理服务器
	fmt.Println("代理服务器已启动: 127.0.0.1:9999")
	http.ListenAndServe("127.0.0.1:9999", nil)
}

// 让Proxy结构体实现 Handle接口
type Proxy struct {
}

func (p *Proxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	fmt.Printf("Received request %s %s %s:\n", req.Method, req.Host, req.RemoteAddr)
	// 1. 代理服务器接收客户端请求，复制，封装成新请求
	outReq := &http.Request{}
	*outReq = *req
	if clientIP, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
		// XFF,用来识别通过HTTP代理或负载均衡方式连接到Web服务器的客户端最原始的IP地址的HTTP请求头字段
		if prior, ok := outReq.Header["X-Forwarded-For"]; ok {
			clientIP = strings.Join(prior, ", ") + ", " + clientIP
		}
		outReq.Header.Set("X-Forwarded-For", clientIP)
	}
	//fmt.Printf("req %s:\n", req.Header.Get("X-Forwarded-For"))
	//fmt.Printf("outReq %s:\n", outReq.Header.Get("X-Forwarded-For"))

	// 2. 发送新请求到下游真实服务器，接收响应
	transport := http.DefaultTransport
	res, err := transport.RoundTrip(outReq)
	if err != nil {
		rw.WriteHeader(http.StatusBadGateway)
		return
	}

	// 3. 处理响应并返回上游客户端
	for key, value := range res.Header {
		for _, v := range value {
			rw.Header().Add(key, v)
		}
	}
	rw.WriteHeader(res.StatusCode)
	io.Copy(rw, res.Body)
	res.Body.Close()
}

func ForwardProxyServer() {
	fmt.Println("正向代理服务器启动 :10000")
	http.Handle("/", &Pxy{})
	http.ListenAndServe("127.0.0.1:10000", nil)
}

// Pxy 定义一个类型，实现 Handler interface
type Pxy struct{}

// ServeHTTP 具体实现方法
func (p *Pxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	fmt.Printf("Received request %s %s %s:\n", req.Method, req.Host, req.RemoteAddr)
	// 1. 代理服务器接收客户端请求，复制，封装成新请求
	outReq := &http.Request{}
	*outReq = *req
	if clientIP, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
		// XFF,用来识别通过HTTP代理或负载均衡方式连接到Web服务器的客户端最原始的IP地址的HTTP请求头字段
		if prior, ok := outReq.Header["X-Forwarded-For"]; ok {
			clientIP = strings.Join(prior, ", ") + ", " + clientIP
		}
		outReq.Header.Set("X-Forwarded-For", clientIP)
	}
	//fmt.Printf("req %s:\n", req.Header.Get("X-Forwarded-For"))
	//fmt.Printf("outReq %s:\n", outReq.Header.Get("X-Forwarded-For"))

	// 2. 发送新请求到下游真实服务器，接收响应
	transport := http.DefaultTransport
	res, err := transport.RoundTrip(outReq)
	if err != nil {
		rw.WriteHeader(http.StatusBadGateway)
		return
	}

	// 3. 处理响应并返回上游客户端
	for key, value := range res.Header {
		for _, v := range value {
			rw.Header().Add(key, v)
		}
	}
	rw.WriteHeader(res.StatusCode)
	io.Copy(rw, res.Body)
	res.Body.Close()
}
