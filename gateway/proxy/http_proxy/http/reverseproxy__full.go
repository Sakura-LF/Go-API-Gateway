package http

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// http反向代理,使用 ReverseProxy 实现
// 功能: URL 重写,更改内容,错误信息回调,连接池

var (
	//代理服务器地址
	proxyAddr = "127.0.0.1:8081"
	// 真实服务器地址
	realServer = "http://127.0.0.1:8001/?a=1"

	// DefaultTransport 连接池配置
	transport = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second, // 连接超时
			KeepAlive: 30 * time.Second, // 长连接超时时间
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,              //最大空闲连接数
		IdleConnTimeout:       90 * time.Second, // 空闲连接超时时间
		TLSHandshakeTimeout:   10 * time.Second, // TLS握手超时时间
		ExpectContinueTimeout: 1 * time.Second,  // 100-countinue 超时时间
	}
)

// ReverseProxySever 使用 ReverseProxy 实现反向代理
func ReverseProxySever() {
	// 解析下游服务器的真实地址
	serverUrl, err := url.Parse(realServer)
	if err != nil {
		log.Fatalln("parse fail,")
	}
	proxy := NewSingleHostReverseProxy(serverUrl)
	log.Println("Starting proxy Server:", proxyAddr)
	// proxy 就相当于 handle
	err = http.ListenAndServe(proxyAddr, proxy)
	if err != nil {
		log.Fatalln(err)
	}
}

// NewSingleHostReverseProxy target代表下游真实服务器地址
func NewSingleHostReverseProxy(target *url.URL) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		// 重写请求函数
		// req 代表请求地址,也就是代理服务器的地址 http://127.0.0.1:8081/sakura
		// target 代表下游真实服务器的地址 // http://127.0.0.1:8001/sakura
		rewriteRequestURL(req, target)
	}

	// 修改返回内容
	modifyResponse := func(res *http.Response) error {
		fmt.Println("enter into modifyResponse function")
		// 升级协议,不需要进行修改
		if res.StatusCode == 101 {
			if strings.Contains(res.Header.Get("Connection"), "Upgrade") {
				return nil
			}
		}
		// 状态嘛等于200 之后在修改
		if res.StatusCode == 200 {
			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				panic(err)
			}
			newBody := []byte("代理服务器返回响应:" + string(resBody))
			// newBody 无法直接赋值给 res.Body
			// res.Body是io.ReadCloser类型,需要实现两个方法
			// 所以直接使用下面的方法
			res.Body = io.NopCloser(bytes.NewBuffer(newBody))
			// Body发生变化后,ContentLength不会变化
			ContentLength := int64(len(newBody))
			res.ContentLength = ContentLength
			// 修改请求头
			res.Header.Set("Content-Length", strconv.FormatInt(ContentLength, 10))
		}
		return nil
	}

	// 错误回调:当后台出现错误响应,会自动调用此函数
	// ModifyResponse 返回 error 也会调用此函数
	errorHandler := func(resp http.ResponseWriter, req *http.Request, err error) {
		fmt.Println("error function")
		http.Error(resp, err.Error(), http.StatusBadGateway)
	}

	// 因为ModifyResponse是reverse的一个属性,所以返回的时候加进去
	return &httputil.ReverseProxy{
		Director:       director,
		ModifyResponse: modifyResponse,
		ErrorHandler:   errorHandler,
		Transport:      transport,
	}
}

func rewriteRequestURL(req *http.Request, target *url.URL) {
	targetQuery := target.RawQuery // 查询参数,?后面的内容
	req.URL.Scheme = target.Scheme // http
	req.URL.Host = target.Host     // 端口号
	// target.path: "" or "/"
	// req.URL.Path: /sakura
	// 合并两个Path, target在前,req在后
	req.URL.Path = joinURLPath(target.Path, req.URL.Path) // /path
	if targetQuery == "" || req.URL.RawQuery == "" {
		req.URL.RawQuery = targetQuery + req.URL.RawQuery
	} else {
		req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
	}
}

// targetPath: "" or "/"
// reqUrlPath: /realserver or ""
func joinURLPath(targetPath, reqUrlPath string) string {
	TString := strings.HasSuffix(targetPath, "/")
	RString := strings.HasPrefix(reqUrlPath, "/")

	switch {
	case TString && RString:
		return targetPath + reqUrlPath[1:] // /sakura -> sakura
	case TString || RString:
		return targetPath + reqUrlPath
	}

	return targetPath + "/" + reqUrlPath
}
