package rpc

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
)

type HelloWorld struct {
}

func (H *HelloWorld) HelloWorld(req string, resp *string) error {
	*resp = req + "你好,Sakura"
	return nil
}

func RpcServer() {
	// 1. 创建服务
	err := rpc.RegisterName("Hello", &HelloWorld{})
	if err != nil {
		log.Println("注册RPC服务失败:", err)
		return
	}

	// 2.设置监听
	listener, err := net.Listen("tcp", "127.0.0.1:8004")
	if err != nil {
		log.Println("监听失败:", err)
	}
	log.Println("监听端口 127.0.0.1:8004")

	// 3.建立连接
	conn, err := listener.Accept()
	if err != nil {
		log.Println("Accept err:", err)
	}
	fmt.Println("连接已建立")
	// 4.绑定服务,将连接和服务进行绑定
	rpc.ServeConn(conn)
}
