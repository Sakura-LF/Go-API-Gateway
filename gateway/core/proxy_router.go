package core

import "C"
import (
	http2 "Go-API-Gateway/gateway/proxy"
	config "Go-API-Gateway/init"
	"Go-API-Gateway/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

var ProxyRouter *gin.Engine

func Proxy() {
	ProxyRouter = gin.New()
	ProxyRouter.Use(gin.Recovery())

	ProxyRouter.GET("/ping", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"msg": "This is Proxy Router",
		})
	})

	ProxyRouter.Run(util.GetHostIp() + ":" + strconv.Itoa(config.Gateway.ProxyPort))
}

func AddProxyRouter(proxyPath string, loadbalance LoadBalance) {

	//根据路径建立代理服务器
	//fmt.Println(proxyPath + servicePath)
	ProxyRouter.GET(proxyPath+"/*name", func(ctx *gin.Context) {
		ctx.Request.URL.Path = ctx.Param("name")
		// 获取路由
		addr, _ := loadbalance.Get(fmt.Sprintf("%s/%v", ctx.Request.Host, ctx.Request.URL))
		addr = "http://" + addr
		parse, err := url.Parse(addr)
		if err != nil {
			log.Fatalln(err)
		}
		proxy := http2.NewMultipleHostsReverseProxy(ctx, parse)

		proxy.ServeHTTP(ctx.Writer, ctx.Request)
	})

	ProxyRouter.POST(proxyPath+"/:name", func(ctx *gin.Context) {
		ctx.Request.URL.Path = ctx.Param("name")
		fmt.Println(ctx.Request.URL.Path)
		// 获取路由
		addr, _ := loadbalance.Get(fmt.Sprintf("%s/%v", ctx.Request.Host, ctx.Request.URL))
		addr = "http://" + addr
		parse, err := url.Parse(addr)
		fmt.Println("parse", parse)
		if err != nil {
			log.Fatalln(err)
		}
		proxy := http2.NewMultipleHostsReverseProxy(ctx, parse)

		proxy.ServeHTTP(ctx.Writer, ctx.Request)
	})
}
