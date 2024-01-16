package consul

import (
	config "Go-API-Gateway/init"
	"encoding/base64"
	"flag"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
	_ "github.com/spf13/viper/remote"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"time"
)

type HttpServer struct {
	server []map[string]int
}

func ServerTest() {
	addr := "127.0.0.1"
	port := 8087
	HTTPServer := HttpServer{server: make([]map[string]int, 10)}
	HTTPServer.server[0][addr] = port
	fmt.Println(HTTPServer)
}

const ConsulAddress = "192.168.74.130:8500"

func ServerRegister() {
	// 处理命令行参数
	Addr := flag.String("addr", "192.168.2.7", "正在监听地址:127.0.0.1")
	Port := flag.Int("port", 8090, "正在监听端口:8090")
	// 拼接地址
	address := fmt.Sprintf("%s:%d", *Addr, *Port)
	flag.Parse()

	// 创建HTTP请求
	service := http.NewServeMux()
	service.HandleFunc("/info", Default)
	service.HandleFunc("/health", Health)

	// 1.定义服务,得到AgentServiceRegistration
	id := uuid.NewString() // 定义注册中心的服务ID
	ProductService := new(api.AgentServiceRegistration)
	ProductService.Name = "Product"
	ProductService.ID = "Product" + id // 使用uuid包创建一个唯一字符串
	ProductService.Address = *Addr
	ProductService.Port = *Port
	ProductService.Tags = []string{"Product"}

	// 2.定义服务健康检查
	ProductService.Checks = api.AgentServiceChecks{
		&api.AgentServiceCheck{
			CheckID:  "product-check-" + id,
			Name:     "Product-Check",
			Interval: "7s",
			Timeout:  "1s",
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
	consulApiConfig.Address = ConsulAddress //地址为consult地址
	consulClient, err := api.NewClient(consulApiConfig)
	if err != nil {
		log.Fatalln(err)
	}

	// 发出 put 注册请求
	if err := consulClient.Agent().ServiceRegister(ProductService); err != nil {
		log.Fatalln(err)
	}
	log.Println("服务注册成功")

	// 启动监听
	fmt.Printf("服务正在监听: %s", address)

	log.Fatalln(http.ListenAndServe(address, service))
}

func Default(writer http.ResponseWriter, request *http.Request) {
	_, err := fmt.Fprintf(writer, "Product Service")
	if err != nil {
		log.Fatal(err)
	}
}

func Health(writer http.ResponseWriter, request *http.Request) {
	log.Println("Health Check")
	_, err := fmt.Fprintf(writer, "Product Service is Health")
	if err != nil {
		log.Fatal(err)
	}
}

func Consulkv() {
	//viper.AddRemoteProvider("consul.yaml", "192.168.74.130:8500", "sakura")
	//viper.SetConfigType("json") // Need to explicitly set this to json
	//err := viper.ReadRemoteConfig()
	//if err != nil {
	//	logs.Println(err)
	//	return
	//}
	//fmt.Println(viper.Get("port")) // 8080
	//fmt.Println(viper.Get("name")) // myhostname.com
	// alternatively, you can create a new viper instance.
	//var runtime_viper = viper.New()
	//
	//runtime_viper.AddRemoteProvider("consul.yaml", "192.168.74.130:8500", "sakura")
	//runtime_viper.SetConfigType("json") // because there is no file extension in a stream of bytes, supported extensions are "json", "toml", "yaml", "yml", "properties", "props", "prop", "env", "dotenv"
	//
	//// read from remote config the first time.
	//err := runtime_viper.ReadRemoteConfig()
	//if err != nil {
	//	logs.Println(err)
	//	return
	//}

	// unmarshal config
	//runtime_viper.Unmarshal(&runtime_conf)

	// open a goroutine to watch remote changes forever
	//go func() {
	//	for {
	//		time.Sleep(time.Second * 5) // delay after each request
	//
	//		// currently, only tested with etcd support
	//		err := runtime_viper.WatchRemoteConfig()
	//		if err != nil {
	//			logs.Printf("unable to read remote config: %v", err)
	//			continue
	//		}
	//
	//		// unmarshal new config into our runtime config struct. you can also use channel
	//		// to implement a signal to notify the system of the changes
	//		//runtime_viper.Unmarshal(&runtime_conf)
	//		fmt.Println(runtime_viper.Get("name"))
	//	}
	//}()
	//select {}
	consulApiConfig := api.DefaultConfig()
	consulApiConfig.Address = ConsulAddress //地址为consult地址
	consulClient, err := api.NewClient(consulApiConfig)
	if err != nil {
		log.Println(err)
		return
	}

	//kvPair := &api.KVPair{Key: "web", Value: make([]byte, 0)}
	//writeOptions := &api.WriteOptions{}
	kvPair, meta, err := consulClient.KV().Get("sakura", &api.QueryOptions{})
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(string(kvPair.Value))
	fmt.Println(meta)
}

//	type mysql struct {
//		Port     int    `yaml:"port"`
//		Name     string `yaml:"name"`
//		Password string `yaml:"password"`
//	}
type Config struct {
	MysqlConfig Mysql `yaml:"mysql"`
}

type Mysql struct {
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

func Kv() {
	var conf Config
	decodeString, err2 := base64.StdEncoding.DecodeString("bXlzcWw6CiAgcG9ydDogMzMwNgogIHVzZXI6IHJvb3QKICBwYXNzd29yZDogc2FrdXJh")
	if err2 != nil {
		log.Println(err2)
	}
	fmt.Println(string(decodeString))

	//file, err2 := os.ReadFile("./mysql.yaml")
	//if err2 != nil {
	//	logs.Fatalln(err2)
	//}
	//fmt.Println(string(file))

	err := yaml.Unmarshal(decodeString, &conf)
	if err != nil {
		panic(err)
	}
	//
	fmt.Println(conf)
}

type Server struct {
	HTTPServer []HTTPServer `yaml:"http_servers" mapstructure:"http_servers"`
}
type HTTPServer struct {
	Addr string `yaml:"addr" mapstructure:"addr"`
	Port int    `yaml:"port" mapstructure:"port"`
}

func TestHTTPServer() {
	//file, err := os.ReadFile("../config/server.yaml")
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(string(file))
	//var server Server
	//err = yaml.Unmarshal(file, &server)
	//if err != nil {
	//	log.Println(err)
	//}
	//fmt.Println(server)
	var server Server
	ServerConfig := config.LoadConfig("server")
	get := ServerConfig.Get("http_servers")
	fmt.Println(get)
	err := ServerConfig.Unmarshal(&server)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(server)
}

type Services struct {
	Name map[string]ServiceObjects
}

type ServiceObjects struct {
	Item map[string]string
}

//type GateWayService struct {
//	Data map[string][]ServiceItems
//}
//
//type ServiceItems struct {
//	Id        string
//	Name      string
//	CreatedAt int64
//	UpdatedAt int64
//	Host      string
//	Protocol  string
//}

func SearchService(name string) {
	consulApiConfig := api.DefaultConfig()
	consulApiConfig.Address = ConsulAddress //地址为consult地址
	consulClient, err := api.NewClient(consulApiConfig)
	if err != nil {
		log.Println(err)
		return
	}
	// 初始化结构体
	//services := Services{Name: make(map[string]ServiceObjects)}
	//objects := ServiceObjects{Item: make(map[string]string)}

	// 查询对应的服务
	filter := "Service==" + name // 拼接filter

	// 每隔10秒查询一次
	for {
		allservices, err := consulClient.Agent().ServicesWithFilter(filter)
		if err != nil {
			log.Fatalln(err)
		}
		for id, service := range allservices {
			host := fmt.Sprintf("%s:%d", service.Address, service.Port)
			fmt.Println("id:", id, " ", host)
			fmt.Println(service.Meta["weight"])

		}
		time.Sleep(time.Second * 10)
	}
}
