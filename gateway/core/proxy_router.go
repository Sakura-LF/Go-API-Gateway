package core

import "C"
import (
	http2 "Go-API-Gateway/gateway/proxy/http_proxy/http"
	config "Go-API-Gateway/init"
	"Go-API-Gateway/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"strconv"
)

var ProxyRouter *gin.Engine

func Proxy() {
	ProxyRouter = gin.Default()

	ProxyRouter.GET("/ping", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"msg": "This is Proxy Router",
		})
	})

	ProxyRouter.Run(util.GetHostIp() + ":" + strconv.Itoa(config.Gateway.ProxyPort))
}

func AddProxyRouter(proxyPath string, urls []*url.URL) {
	//根据路径建立代理服务器
	//fmt.Println(proxyPath + servicePath)
	ProxyRouter.GET(proxyPath+"/:name", func(ctx *gin.Context) {
		fmt.Println(ctx.Param("name"))
		ctx.Request.URL.Path = ctx.Param("name")

		proxy := http2.NewMultipleHostsReverseProxy(ctx, urls)

		proxy.ServeHTTP(ctx.Writer, ctx.Request)
	})
}
