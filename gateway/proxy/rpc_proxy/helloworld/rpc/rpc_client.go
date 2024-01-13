package rpc

import (
	"fmt"
	"log"
	"net/rpc"
)

func RpcClient() {

	// 1.用rpc链接服务器 --Dial(拨号)
	conn, err := rpc.Dial("tcp", "127.0.0.1:8004")
	if err != nil {
		log.Println("Dial err:", err)
	}
	defer conn.Close()

	// 2.调用远程函数
	var response string // 接收返回值
	// 指定调用的远程函数
	err = conn.Call("Hello.HelloWorld", "Sakura", &response)
	if err != nil {
		log.Fatalln("Call err:", err)
	}
	fmt.Println(response)
}
