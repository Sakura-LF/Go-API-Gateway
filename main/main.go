package main

import (
	"Go-API-Gateway/gateway/core"
	config "Go-API-Gateway/init"
	"github.com/rs/zerolog/log"
)

func init() {
	// 读取所有配置文件,初始化所有配置结构体
	config.InitConfig()
}

func main() {
	go core.Api()
	go core.Proxy()
	// 并且整合到负载均衡

	// 新建一个发布者
	core.NewConcreteSubject()

	log.Debug().Msg("API 网关启动!")
	log.Error().Msg("API 网关启动!")
	log.Warn().Msg("API 网关启动!")
	log.Info().Msg("API 网关启动!")
	select {}
}
