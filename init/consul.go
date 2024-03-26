package init

import (
	"github.com/hashicorp/consul/api"
	"github.com/rs/zerolog/log"
)

var ConsulClient *api.Client

func InitConsul() {
	consulApiConfig := api.DefaultConfig()
	consulApiConfig.Address = ConsulConfig.GetString("consul.addr") //地址为consult地址
	consulClient, err := api.NewClient(consulApiConfig)
	if err != nil {
		log.Error().Err(err).Msg("consulClient init fail:")
	}
	ConsulClient = consulClient
}
