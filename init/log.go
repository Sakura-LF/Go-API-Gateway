package init

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"strings"
	"time"
)

// LogConfig 日志配置
type LogConfig struct {
	ZeroLogConfig ZeroLogConfig `yaml:"ZeroLogConfig"`
	LogRotate     LogRotate     `yaml:"LogRotate"`
}

// ZeroLogConfig
type ZeroLogConfig struct {
	Level   string `yaml:"Level"`
	Pattern string `yaml:"Pattern"`
	OutPut  string `yaml:"OutPut"`
}

// LogRotate  日志轮换(分割)配置
type LogRotate struct {
	Filename   string `yaml:"Filename"`
	MaxSize    int    `yaml:"MaxSize"`    // megabytes，M 为单位，达到这个设置数后就进行日志切割
	MaxBackups int    `yaml:"MaxBackups"` // 保留旧文件最大份数
	MaxAge     int    `yaml:"MaxAge"`     // days ， 旧文件最大保存天数
	Compress   bool   `yaml:"Compress"`   // 是否开启压缩,默认关闭
}

// LogInit 完成Zero 日志的初始化
func LogInit() {
	var logConfig LogConfig
	config := LoadConfig("log")
	// 反序列化到结构体
	err := config.Unmarshal(&logConfig)
	if err != nil {
		panic(err)
	}
	// 验证是序列化成功
	//fmt.Println(logConfig.ZeroLogConfig)
	//fmt.Println(logConfig.LogRotate)
	//fmt.Println(strings.Join([]string{logConfig.ZeroLogConfig.OutPut, logConfig.LogRotate.Filename}, "/"))

	// 设置日志等级
	switch logConfig.ZeroLogConfig.Level {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	}

	// 日志切割
	logRotate := &lumberjack.Logger{
		Filename:   strings.Join([]string{logConfig.ZeroLogConfig.OutPut, logConfig.LogRotate.Filename}, "/"), // 文件位置
		MaxSize:    1,                                                                                         // megabytes，M 为单位，达到这个设置数后就进行日志切割
		MaxBackups: 3,                                                                                         // 保留旧文件最大份数
		MaxAge:     7,                                                                                         //days ， 旧文件最大保存天数
		Compress:   true,                                                                                      // disabled by default，是否压缩日志归档，默认不压缩
	}
	// 调整日志时间格式
	zerolog.TimeFieldFormat = time.StampMilli
	// 开启调用位置打印
	log.Logger = log.With().Caller().Logger()

	if logConfig.ZeroLogConfig.Pattern == "development" {
		// 控制台输出的输出器
		consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.StampMilli}
		multi := zerolog.MultiLevelWriter(consoleWriter, logRotate)
		log.Logger = log.Output(multi)
	} else if logConfig.ZeroLogConfig.Pattern == "production" {
		log.Logger = log.Output(logRotate)
	} else {
		panic("log pattern Error")
	}
}
