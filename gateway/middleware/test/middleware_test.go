package test

import (
	"Go-API-Gateway/gateway/middleware/router"
	"Go-API-Gateway/gateway/middleware/timerate"
	http2 "Go-API-Gateway/gateway/proxy"
	"context"
	"fmt"
	"golang.org/x/time/rate"
	"log"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestSlice(t *testing.T) {
	// 代理服务器地址
	var addr = "127.0.0.1:8006"

	// 1.创建服务器
	// 一个路由器,包含多个路由
	// 每个路由器都可以又多个处理器(回调函数)
	sliceRouter := router.NewSliceRouter()

	// 2.构建URI路由中间件,注册请求URI
	routeRoot := sliceRouter.Group("/")

	// 3.为路由绑定处理函数
	routeRoot.Use(handle, func(c *router.SliceRouteContext) {
		fmt.Println("reverse proxy")
		reverseProxy(c.Ctx).ServeHTTP(c.RW, c.Req)
	})

	// 将路由器作为http服务的处理器
	//  封装 sliceRouter 作为http服务的处理器
	routerHandler := router.NewSliceRouterHandler(nil, sliceRouter)
	log.Println("Starting httpserver at " + addr)
	http.ListenAndServe(addr, routerHandler)

}

func handle(c *router.SliceRouteContext) {
	log.Println("trace...in")
	//c.Next()
	log.Println("trace...out")
}

func reverseProxy(c context.Context) http.Handler {
	rs1 := "http://127.0.0.1:8001/"
	url1, err1 := url.Parse(rs1)
	if err1 != nil {
		log.Println(err1)
	}

	rs2 := "http://127.0.0.1:8002/"
	url2, err2 := url.Parse(rs2)
	if err2 != nil {
		log.Println(err2)
	}

	urls := []*url.URL{url1, url2}
	return http2.NewMultipleHostsReverseProxy(c, urls)
}

// 使用 golang/org/x/time/rate 包实现一个限流器
// 实现步骤
// 1.构建一个限速器
// 2.获取Token
//
//	Wait 阻塞,知道获取token
//	Reserve 预约,等待指定时间,再获取token
//	Allow 返回bool值,判断当前是否可以获取token
func TestTimeRate(t *testing.T) {
	// 1.构建限速器
	// 第一个参数 r : 每秒产生的token数量
	// 第二个参数 l : 对到的token数量(令牌桶的容量)
	limiter := rate.NewLimiter(1, 10)

	// 2.获取Token,三种方式
	for i := 0; i < 200; i++ {
		t.Log("before wait:", i)
		// 阻塞等待直到获取一个token
		ctx, _ := context.WithTimeout(context.Background(), time.Second*2)
		if err := limiter.Wait(ctx); err != nil {
			t.Log("timeout:", err)
		}
		t.Log("after wait:")

		// 返回预计多久才能有新的token,可以等待指定时间获取
		reserve := limiter.Reserve()
		// 检查限速器能否在等待时间内返回token
		if !reserve.OK() {
			return
		}
		t.Log("reserve time:", reserve.Delay())
		time.Sleep(reserve.Delay())

		// 判断当前是否可以取到token
		// 若返回true,则已经获取到token
		t.Log("Allow:", limiter.Allow())

		time.Sleep(time.Millisecond * 300)
	}

}

func TestRateLimiter2(t *testing.T) {
	customHandler := func(c *router.SliceRouteContext) http.Handler {
		rs1 := "http://127.0.0.1:8001/"
		url1, err1 := url.Parse(rs1)
		if err1 != nil {
			log.Println(err1)
		}

		rs2 := "http://127.0.0.1:8002/"
		url2, err2 := url.Parse(rs2)
		if err2 != nil {
			log.Println(err2)
		}

		urls := []*url.URL{url1, url2}
		return http2.NewMultipleHostsReverseProxy(c.Ctx, urls)
	}
	// 代理服务器
	addr := "127.0.0.1:8006"
	log.Println("Starting http Server at :", addr)
	r := router.NewSliceRouter()
	r.Group("/").Use(timerate.RateLimiter())
	routerHandler := router.NewSliceRouterHandler(customHandler, r)
	log.Fatalln(http.ListenAndServe(addr, routerHandler))
}
