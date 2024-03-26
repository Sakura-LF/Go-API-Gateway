package init

import (
	"Go-API-Gateway/util"
	"errors"
	"github.com/spf13/viper"
)

type GateWay struct {
	ProxyPort int `mapstructure:"proxy_port"`
	ApiPorts  int `mapstructure:"api_port"`
}

var (
	Gateway      *GateWay
	ConsulConfig *viper.Viper
)

func InitConfig() {
	LogInit()
	ConsulConfig = LoadConfig("consul")
	UnmarshalStruct("gateway")
	InitConsul()
}

func UnmarshalStruct(filename string) {
	GateWayConfig := LoadConfig(filename)
	err := GateWayConfig.Unmarshal(&Gateway)
	if err != nil {
		panic(err)
	}
}

func LoadConfig(filename string) *viper.Viper {
	config := viper.New()
	rootPath := util.GetRootPath()
	config.AddConfigPath(rootPath + "/config")
	config.SetConfigName(filename)
	err := config.ReadInConfig()
	if err != nil {
		// 如果需要对配置文件不存在错误，做特殊处理，使用：
		if errors.As(err, &viper.ConfigFileNotFoundError{}) {
			//config file not found; ignore error if desired
			panic("config file was not found: " + filename)
		} else {
			panic(err)
		}
	}
	return config
}
