package core

import (
	config "Go-API-Gateway/init"
	"Go-API-Gateway/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"time"
)

func Api() {
	api := gin.Default()

	api.GET("/ping", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"msg": "This is API Router",
		})
	})

	// 添加服务
	api.POST("/services", func(ctx *gin.Context) {
		//name := ctx.PostForm("name")
		name, _ := ctx.GetPostForm("name")
		//wg := &sync.WaitGroup{}
		go LoadService(name)
		//从consul中查出服务名称
		//filter := "Service==" + name
		//allservices, err := config.ConsulClient.Agent().ServicesWithFilter(filter)
		//if err != nil {
		//	log.Fatalln(err)
		//}

		//for id, service := range allservices {
		//	host := fmt.Sprintf("%s:%d", service.Address, service.Port)
		//	serviceItems = ServiceItems{
		//		Id:        id,
		//		Name:      name,
		//		Host:      host,
		//		CreatedAt: time.Now().Unix(),
		//		UpdatedAt: time.Now().Unix(),
		//		Protocol:  "http",
		//	}
		//	Service.Data[name] = append(Service.Data[name], serviceItems)
		//}

		//fmt.Println(serviceItems)

		//Service.data = append(Service.data[name], serviceItems)
		ctx.JSON(http.StatusOK, Service.Data[name])
	})

	api.GET("/services", func(ctx *gin.Context) {

		ctx.JSON(http.StatusOK, Service.Data)
	})

	// 路由接口
	// 添加路由
	api.POST("/routers", func(ctx *gin.Context) {
		// 指定要把那个或者哪些路径添加到对应的服务上
		servicename, _ := ctx.GetPostForm("servicename")
		//fmt.Println(service)
		// 指定路径
		path, _ := ctx.GetPostForm("paths")
		loadbalanceType, err := ctx.GetPostForm("loadbalance")
		if !err {
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
		fmt.Println(loadblance.Next())
		fmt.Println(loadblance.Next())
		//fmt.Println(urls)
		// 建立代理服务器
		AddProxyRouter(path, loadblance)

		// todo 对应的负载均衡算法选择

	})

	api.GET("/routers", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, Router.data)
	})

	api.Run(util.GetHostIp() + ":" + strconv.Itoa(config.Gateway.ApiPorts))
}

// // TODO
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
