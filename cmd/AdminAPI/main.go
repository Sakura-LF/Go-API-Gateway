package main

import (
	"Go-API-Gateway/gateway/core"
	config "Go-API-Gateway/init"
)

func init() {
	// 读取所有配置文件,初始化所有配置结构体
	config.InitConfig()
}

func main() {
	core.Api()
	core.NewConcreteSubject()
}
