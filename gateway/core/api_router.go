package core

import (
	config "Go-API-Gateway/init"
	"Go-API-Gateway/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

func Api() {
	api := gin.New()
	api.Use(gin.Recovery())
	//api.Use(gin.Logger())

	// 检测 admin api服务
	api.GET("/ping", adminAPiPing)
	// 添加服务
	api.POST("/services", addService)
	// 获取所有服务信息
	api.GET("/services", getService)

	// 路由接口
	// 添加路由
	api.POST("/routers", addRouter)
	// 获取路由
	api.GET("/routers", getRouter)

	// 插件管理
	api.POST("/plugin/service/ratelimit", addPlugin)
	//api.POST("/plugin/service/")

	api.Run(util.GetHostIp() + ":" + strconv.Itoa(config.Gateway.ApiPorts))
}

func addPlugin(ctx *gin.Context) {
	servicename, b := ctx.GetPostForm("servicename")
	if !b {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "servicename 必须指定",
		})
	}
	second, b := ctx.GetPostForm("second")
	if !b {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "second 必须指定",
		})
	}
	burst, b := ctx.GetPostForm("burst")
	if !b {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "second 必须指定",
		})
	}
	fmt.Println(second)
	fmt.Println(burst)
	routerItems := Router.data[servicename]
	for key, value := range routerItems {
		fmt.Println(key, " ", value)
	}
	fmt.Println(routerItems)
}

func adminAPiPing(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"msg": "This is AdminAPI Router",
	})
}

func addService(ctx *gin.Context) {
	//name := ctx.PostForm("name")
	name, _ := ctx.GetPostForm("name")
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go LoadService(name, wg)
	wg.Wait()
	ctx.JSON(http.StatusOK, Service.Data[name])
	log.Info().Str("ServiceName", name).Msg("服务添加成功")
}

func getService(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, Service.Data)
	log.Info().Msg("获取 API 网关所有服务")
}

func addRouter(ctx *gin.Context) {
	// 指定要把那个或者哪些路径添加到对应的服务上
	servicename, _ := ctx.GetPostForm("servicename")
	//fmt.Println(service)
	// 指定路径
	path, _ := ctx.GetPostForm("paths")
	loadbalanceType, err := ctx.GetPostForm("loadbalance")
	if !err {
		log.Error().Msg("没有选择对应的负载均衡器,路由添加失败")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "loadbalance must be choose form Random or RoundRobin or WeightRoundRobin or ConsistenHash ",
		})
		return
	}

	items := RouterItems{
		Id:        uuid.NewString(),
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
		Path:      nil,
		Protocol:  "",
	}
	items.Path = append(items.Path, path)

	Router.data[servicename] = append(Router.data[servicename], items)

	// 响应JSON
	ctx.JSON(http.StatusOK, Router.data)

	// 选择负载均衡算法
	loadblance := ChooseLoadBalance(loadbalanceType, ctx)
	// fmt.Println(loadblance.Next())
	// fmt.Println(loadblance.Next())
	// fmt.Println(urls)
	// 建立代理服务器
	AddProxyRouter(path, loadblance)
	log.Info().Str("Route", path).Str("ServiceName", servicename).Msg("路由添加成功")
}

func getRouter(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, Router.data)
	log.Info().Msg("获取所有路由")
}

func ChooseLoadBalance(loadbalanceType string, ctx *gin.Context) LoadBalance {
	// 判断负载均衡类型
	if loadbalanceType == "RoundRobin" {
		factory := LoadBalanceFactory(RoundRobin, ConcreteSubjectConsul)
		return factory
	} else if loadbalanceType == "Random" {
		factory := LoadBalanceFactory(Random, ConcreteSubjectConsul)
		return factory
	} else if loadbalanceType == "WeightRoundRobin" {
		factory := LoadBalanceFactory(RoundRobin, ConcreteSubjectConsul)
		return factory
	} else if loadbalanceType == "ConsistenHash" {
		factory := LoadBalanceFactory(RoundRobin, ConcreteSubjectConsul)
		return factory
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "loadbalance must be choose form Random or RoundRobin or WeightRoundRobin or ConsistenHash ",
		})
		return nil
	}
}
