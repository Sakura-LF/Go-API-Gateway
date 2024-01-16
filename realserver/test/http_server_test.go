package test

import (
	"Go-API-Gateway/realserver"
	"testing"
)

func TestRealHTTPServer(t *testing.T) {
	// 启动三个真实服务,并且将服务注册到服务发现中心,开启健康检查
	realserver.RealHTTPServer("192.168.2.7", 9000, "40")
	realserver.RealHTTPServer("192.168.2.7", 9001, "20")
	realserver.RealHTTPServer("192.168.2.7", 9002, "10")
	realserver.RealHTTPServer("192.168.2.7", 9003, "5")
	select {}
}

func TestRealHTTPServer2(t *testing.T) {
	// 启动三个真实服务,并且将服务注册到服务发现中心,开启健康检查
	realserver.RealHTTPServer("192.168.2.7", 9004, "100")

	select {}
}

func TestRealHTTPServer3(t *testing.T) {
	// 启动三个真实服务,并且将服务注册到服务发现中心,开启健康检查
	realserver.RealHTTPServer("192.168.2.7", 9005, "200")

	select {}
}
