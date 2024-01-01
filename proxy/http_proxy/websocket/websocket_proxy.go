package websocket

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var (
	// 代理服务器地址
	proxyServer = "127.0.0.1:8082"
	// 真实websocket服务器地址
	websocketServer = "http://127.0.0.1:8002"
)

func WebSocketProxy() {
	url, err := url.Parse(websocketServer)
	if err != nil {
		log.Println(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(url)
	log.Println("WebSocket 代理启动, 按CTRL+C退出")
	http.ListenAndServe(proxyServer, proxy)
}
