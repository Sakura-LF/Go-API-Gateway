package tcp_proxy_test

import (
	tcp_proxy2 "Go-API-Gateway/proxy/tcp_proxy"
	"context"
	"fmt"
	"net"
	"testing"
)

func TestRunTCPServer(t *testing.T) {
	// 启动TCP服务器
	go func() {
		addr := "192.168.2.7:8001"
		// 1.启动一个TCP Server
		tcpServer := &tcp_proxy2.TCPServer{
			Addr:    addr,
			Handler: &MyHander{},
		}
		// 2. 启动TCP服务器
		fmt.Println("Starting TCP Server at:", addr)
		tcpServer.ListenAndServe()
	}()

	// 启动TCP代理服务器
	go func() {
		tcpServerAddr := "192.168.2.7:8001"
		// 1.创建TCP代理实例
		tcpProxy := tcp_proxy2.NewSingleHostReverseProxy(tcpServerAddr)
		// 2. 启动TCP代理服务器
		tcpProxyAddr := "192.168.2.7:8081"
		fmt.Println("Starting TCP Proxy Server at:", tcpProxyAddr)
		tcp_proxy2.ListenAndServe(tcpProxyAddr, tcpProxy)
	}()
	select {}
}

type MyHander struct {
}

func (handler *MyHander) ServeTCP(ctx context.Context, conn net.Conn) {
	conn.Write([]byte("Sakura ! ! ! \n"))
}
