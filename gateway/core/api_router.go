package core

import (
	config "Go-API-Gateway/init"
	"Go-API-Gateway/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func Api() {
	api := gin.Default()

	var serviceItems ServiceItems

	api.GET("/ping", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"msg": "This is API Router",
		})
	})

	// 添加服务
	api.POST("/services", func(ctx *gin.Context) {
		//name := ctx.PostForm("name")
		name, _ := ctx.GetPostForm("name")

		//从consul中查出服务名称
		filter := "Service==" + name
		allservices, err := config.ConsulClient.Agent().ServicesWithFilter(filter)
		if err != nil {
			log.Fatalln(err)
		}

		for id, service := range allservices {
			host := fmt.Sprintf("%s:%d", service.Address, service.Port)
			serviceItems = ServiceItems{
				Id:        id,
				Name:      name,
				Host:      host,
				CreatedAt: time.Now().Unix(),
				UpdatedAt: time.Now().Unix(),
				Protocol:  "http",
			}
			Service.data[name] = append(Service.data[name], serviceItems)
		}

		//fmt.Println(serviceItems)

		//Service.data = append(Service.data[name], serviceItems)
		ctx.JSON(http.StatusOK, Service.data[name])
	})

	api.GET("/services", func(ctx *gin.Context) {

		ctx.JSON(http.StatusOK, Service.data)
	})

	// 路由接口
	// 添加路由
	api.POST("/routers", func(ctx *gin.Context) {
		// 指定要把那个或者哪些路径添加到对应的服务上
		servicename, _ := ctx.GetPostForm("servicename")
		//fmt.Println(service)
		// 指定路径
		path, _ := ctx.GetPostForm("paths")

		items := RouterItems{
			Id:        uuid.NewString(),
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
			Path:      nil,
			Protocol:  "",
		}
		items.Path = append(items.Path, path)

		Router.data[servicename] = append(Router.data[servicename], items)

		ctx.JSON(http.StatusOK, Router.data)

		// todo 创建代理和路由进行绑定
		urls := make([]*url.URL, 0)
		// 首先遍历服务地址
		services := Service.data[servicename]
		for _, value := range services {
			host, err := url.Parse(fmt.Sprintf("http://%s", value.Host))
			if err != nil {
				log.Println("parse failed")
			}
			//fmt.Println(host)
			urls = append(urls, host)
			// 代理服务地址
		}
		//fmt.Println(urls)
		// 建立代理服务器
		AddProxyRouter(path, urls)

	})

	api.GET("/routers", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, Router.data)
	})

	api.Run(util.GetHostIp() + ":" + strconv.Itoa(config.Gateway.ApiPorts))
}
