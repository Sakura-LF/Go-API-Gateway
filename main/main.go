package main

import (
	"Go-API-Gateway/gateway/core"
	config "Go-API-Gateway/init"
	"fmt"
)

func init() {
	// 读取所有配置文件,初始化所有配置结构体
	config.InitConfig()
	config.ZapInit()
}

func main() {
	//ip := getHostIp()
	//fmt.Println(ip)
	fmt.Println(config.Gateway)

	go core.Api()
	go core.Proxy()
	// 并且整合到负载均衡

	// 新建一个发布者
	core.NewConcreteSubject()

	select {}
}
