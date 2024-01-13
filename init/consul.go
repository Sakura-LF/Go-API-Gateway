package init

import (
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
)

var ConsulClient *api.Client

func InitConsul() {
	consulApiConfig := api.DefaultConfig()
	consulApiConfig.Address = ConsulConfig.GetString("consul.addr") //地址为consult地址
	consulClient, err := api.NewClient(consulApiConfig)
	if err != nil {
		Logger.Info("consulClient init fail:", zap.Error(err))
	}
	ConsulClient = consulClient
}
