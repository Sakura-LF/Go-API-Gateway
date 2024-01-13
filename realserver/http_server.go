package realserver

import (
	config "Go-API-Gateway/init"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
	"log"
	"net/http"
	"strconv"
	"time"
)

// consul地址
var consulAddress = "192.168.74.130:8500"

func ByConfigRunHTTPServer() {

}

func RealHTTPServer(addr string, port int) {
	server := &RealServer{Addr: addr}
	// 向consul注册服务
	//Register(addr, port)
	address := fmt.Sprintf("%s:%d", addr, port)

	// 1.定义服务,得到AgentServiceRegistration
	id := strconv.Itoa(port) // 定义注册中心的服务ID
	ProductService := new(api.AgentServiceRegistration)
	ProductService.Name = "Test"
	ProductService.ID = "Test" + id // 使用uuid包创建一个唯一字符串
	ProductService.Address = addr
	ProductService.Port = port
	ProductService.Tags = []string{"Test"}

	// 2.定义服务健康检查
	ProductService.Checks = api.AgentServiceChecks{
		&api.AgentServiceCheck{
			CheckID:  "Test-check-" + id,
			Name:     "Test-Check",
			Interval: "10s",
			Timeout:  "2s",
			// 请求地址一定不要拼接错
			HTTP:                           fmt.Sprintf("http://%s/ping", address),
			Method:                         "GET",
			SuccessBeforePassing:           0,
			FailuresBeforeWarning:          0,
			FailuresBeforeCritical:         0,
			DeregisterCriticalServiceAfter: "",
		},
	}

	// 2.获取consul连接
	// 2.1 配置consul 服务器地址
	config.InitConfig()
	config.InitConsul()
	fmt.Println(config.ConsulConfig.Get("consul.addr"))

	// 发出 put 注册请求
	if err := config.ConsulClient.Agent().ServiceRegister(ProductService); err != nil {
		log.Fatalln("register fail", err)
	}
	log.Println("服务注册成功")

	// 启动http服务
	go func() {
		server.Run(addr, port)
	}()
}

func (r *RealServer) Run(addr string, port int) {
	address := fmt.Sprintf("%s:%d", addr, port)

	mux := http.NewServeMux()
	mux.HandleFunc("/info", r.InfoHandler)
	mux.HandleFunc("/index", r.IndexHandler)
	mux.HandleFunc("/ping", r.Health)

	server := http.Server{
		Addr:         address,
		Handler:      mux,
		WriteTimeout: time.Second * 3,
	}

	// 开启协程启动监听
	fmt.Println("http服务器已启动:", address)
	server.ListenAndServe()
}

type RealServer struct {
	Addr string
}

func (r *RealServer) InfoHandler(w http.ResponseWriter, req *http.Request) {
	// 拼接真实服务器地址
	URL := fmt.Sprintf("Server Info,address: http://%s%s", req.Host, req.URL.Path)
	w.Write([]byte(URL))
}

func (r *RealServer) IndexHandler(w http.ResponseWriter, req *http.Request) {
	// 拼接真实服务器地址
	URL := fmt.Sprintf("Server Index,address: http://%s%s", req.Host, req.URL.Path)
	w.Write([]byte(URL))
}

func (r *RealServer) Health(writer http.ResponseWriter, request *http.Request) {
	log.Println(request.Host, " Health Check")
	_, err := fmt.Fprintf(writer, "Service is Health")
	if err != nil {
		log.Fatal(err)
	}
}

func Register(addr string, port int) {
	address := fmt.Sprintf("%s:%d", addr, port)

	// 1.定义服务,得到AgentServiceRegistration
	id := uuid.NewString() // 定义注册中心的服务ID
	ProductService := new(api.AgentServiceRegistration)
	ProductService.Name = "Test"
	ProductService.ID = "Test" + id // 使用uuid包创建一个唯一字符串
	ProductService.Address = addr
	ProductService.Port = port
	ProductService.Tags = []string{"Test"}

	// 2.定义服务健康检查
	ProductService.Checks = api.AgentServiceChecks{
		&api.AgentServiceCheck{
			CheckID:  "Test-check-" + id,
			Name:     "Test-Check",
			Interval: "10s",
			Timeout:  "2s",
			// 请求地址一定不要拼接错
			HTTP:                           fmt.Sprintf("http://%s/health", address),
			Method:                         "GET",
			SuccessBeforePassing:           0,
			FailuresBeforeWarning:          0,
			FailuresBeforeCritical:         0,
			DeregisterCriticalServiceAfter: "",
		},
	}

	// 2.注册服务
	// 2.1 配置consul 服务器地址
	consulApiConfig := api.DefaultConfig()
	consulApiConfig.Address = consulAddress //地址为consult地址
	consulClient, err := api.NewClient(consulApiConfig)
	if err != nil {
		log.Fatalln(err)
	}

	// 发出 put 注册请求
	if err := consulClient.Agent().ServiceRegister(ProductService); err != nil {
		log.Fatalln(err)
	}
	log.Println("服务注册成功")
}
