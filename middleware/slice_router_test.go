package middleware

import (
	"fmt"
	"log"
	"net/http"
	"testing"
)

func TestSlice(t *testing.T) {
	// 代理服务器地址
	var addr = "127.0.0.1:8006"

	// 1.创建服务器
	// 一个路由器,包含多个路由
	// 每个路由器都可以又多个处理器(回调函数)
	sliceRouter := NewSliceRouter()

	// 2.构建URI路由中间件,注册请求URI
	routeRoot := sliceRouter.Group("/")

	// 3.为路由绑定处理函数
	routeRoot.Use(handle, func(c *SliceRouteContext) {
		fmt.Println("reverse proxy")
	})

	// 将路由器作为http服务的处理器
	// TODO 封装 sliceRouter 作为http服务的处理器
	routerHandler := NewSliceRouterHandler(nil, sliceRouter)
	http.ListenAndServe(addr, routerHandler)
}

func handle(c *SliceRouteContext) {
	log.Println("trace...in")
	c.Next()
	log.Println("trace...out")
}
