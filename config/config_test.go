package config

import (
	config "Go-API-Gateway/init"
	"fmt"
	"github.com/spf13/viper"
	"testing"
)

type Server struct {
	HTTPServer []HTTPServer
}
type HTTPServer struct {
	addr map[string]int
}

func TestConfig(t *testing.T) {
	viper.AddConfigPath(".")
	viper.SetConfigName("server")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	slice := viper.Get("http_servers")
	fmt.Println(slice)
	var AllServer Server
	//var Httpserver HTTPServer
	err = viper.UnmarshalKey("http_servers", &AllServer.HTTPServer)
	if err != nil {
		panic(err)
	}
	//get
	//if err != nil {
	//	panic(err)
	//}
	fmt.Println(AllServer)
}

type GateWay struct {
	ProxyPort int `mapstructure:"proxy_port"`
	ApiPorts  int `mapstructure:"api_port"`
}

func TestUnm(t *testing.T) {
	var gateway GateWay
	//file, _ := os.ReadFile("./gateway.yaml")
	//fmt.Println(string(file))
	//err := yaml.Unmarshal(file, &gateWay)
	//if err != nil {
	//	panic(err)
	//}
	GateWayConfig := config.LoadConfig("gateway")
	fmt.Println(GateWayConfig.Get("gateway.proxy_port"))
	err := GateWayConfig.UnmarshalKey("gateway.proxy_port", &gateway.ProxyPort)
	if err != nil {
		panic(err)
	}
	fmt.Println(gateway)
}
