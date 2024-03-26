package core

import (
	config "Go-API-Gateway/init"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"
)

var Service *GateWayService
var Router *GateRouter

func init() {
	Service = &GateWayService{Data: make(map[string][]ServiceItems, 10)}
	//Items = make([]ServiceItems, 0)
	Router = &GateRouter{data: make(map[string][]RouterItems, 10)}
}

// type Get interface {
// }
func GetService() *GateWayService {
	return Service
}

type GateWayService struct {
	Data map[string][]ServiceItems
}

type ServiceItems struct {
	Id        string
	Name      string
	CreatedAt int64
	UpdatedAt int64
	Host      string
	Weight    int
	Protocol  string
}

type GateRouter struct {
	data map[string][]RouterItems
}

type RouterItems struct {
	Id        string
	CreatedAt int64
	UpdatedAt int64
	Path      []string
	Protocol  string
	Plugin    []PluginConfig
}

type PluginConfig struct {
	RateLimiter RateLimiter
}

// RateLimiter limiter := rate.NewLimiter(rate.Every(time.Second), 1)
type RateLimiter struct {
	Name   string
	Second int
	Burst  int
}

// LoadService 从consul中将服务绑定到Service中
func LoadService(name string, wg *sync.WaitGroup) {
	var serviceItems ServiceItems

	// 查询对应的服务
	filter := "Service==" + name // 拼接filter
	servicesLength := 0
	// 每隔10秒查询一次
	for {
		allservices, err := config.ConsulClient.Agent().ServicesWithFilter(filter)
		if err != nil {
			log.Fatalln(err)
		}
		servicesLength = len(allservices)
		if len(Service.Data[name]) == 0 {
			for id, service := range allservices {
				host := fmt.Sprintf("%s:%d", service.Address, service.Port)
				metaWeight := service.Meta["weight"]
				weight, _ := strconv.Atoi(metaWeight)
				//fmt.Println(host)
				serviceItems = ServiceItems{
					Id:        id,
					Name:      name,
					Host:      host,
					Weight:    weight,
					CreatedAt: time.Now().Unix(),
					UpdatedAt: time.Now().Unix(),
					Protocol:  "http",
				}
				Service.Data[name] = append(Service.Data[name], serviceItems)
			}
			wg.Done()
			//fmt.Println("-----------")
			//fmt.Println(Service)
			time.Sleep(time.Second * 10)
			continue
		} else if len(Service.Data[name]) != servicesLength {
			Service = &GateWayService{Data: make(map[string][]ServiceItems, 10)}
			continue
		} else if len(Service.Data[name]) == servicesLength {
			//fmt.Println("len(core.Service.Data[name]):", len(Service.Data[name]), "servicesLength:", servicesLength)
			time.Sleep(time.Second * 10)
			continue
		}
	}
}
