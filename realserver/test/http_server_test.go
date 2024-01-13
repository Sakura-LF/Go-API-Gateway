package test

import (
	"Go-API-Gateway/realserver"
	"testing"
)

func TestRealHTTPServer(t *testing.T) {
	// 启动三个真实服务,并且将服务注册到服务发现中心,开启健康检查
	realserver.RealHTTPServer("192.168.2.7", 9000)
	realserver.RealHTTPServer("192.168.2.7", 9001)
	realserver.RealHTTPServer("192.168.2.7", 9002)
	realserver.RealHTTPServer("192.168.2.7", 9003)
	select {}
}
