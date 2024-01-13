package test

import (
	conf "Go-API-Gateway/init"
	"fmt"
	"log"
	"testing"
)

type Server struct {
	HTTPServer []HTTPServer `yaml:"http_servers"`
}
type HTTPServer struct {
	Addr string `yaml:"addr"`
	Port int    `yaml:"port"`
}

func TestConfig(t *testing.T) {
	config := conf.LoadConfig("server")
	//get := config.Get("http_servers")
	//fmt.Println(get)
	//conf.UnmarshalConfig()
	//fmt.Println(conf.ServerConf)
	var server Server
	err := config.Unmarshal(&server)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(server)
}
