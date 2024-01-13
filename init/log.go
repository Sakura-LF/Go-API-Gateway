package init

import (
	"Go-API-Gateway/util"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

var (
	Logger *zap.Logger
)

func ZapInit() {
	level := zap.DebugLevel
	levelstring := LogConfig.GetString("zap.level")

	switch levelstring {
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	}
	path := util.GetRootPath()
	fmt.Println(path)
	OutPutSlice := LogConfig.GetStringSlice("zap.outputPaths")
	OutPutSlice[1] = path + OutPutSlice[1]
	ErrOutPutSlice := LogConfig.GetStringSlice("zap.errorOutputPaths")
	ErrOutPutSlice[1] = path + ErrOutPutSlice[1]

	zapConfig := zap.Config{
		Level:       zap.NewAtomicLevelAt(level),
		Development: true,
		Encoding:    LogConfig.GetString("zap.encoding"),
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:     "time",
			LevelKey:    "level",
			NameKey:     "logger",
			CallerKey:   "caller", // 记录日志调用位置
			FunctionKey: zapcore.OmitKey,
			MessageKey:  "message",
			//StacktraceKey: "Stack",
			LineEnding:  zapcore.DefaultLineEnding,
			EncodeLevel: zapcore.LowercaseLevelEncoder,
			EncodeTime: func(time time.Time, encoder zapcore.PrimitiveArrayEncoder) {
				encoder.AppendString(time.Local().Format("2006-04-02 15:04:05"))
			},
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      OutPutSlice,
		ErrorOutputPaths: ErrOutPutSlice,
	}
	//config := zap.NewDevelopmentConfig()
	Logger = zap.Must(zapConfig.Build())
}
