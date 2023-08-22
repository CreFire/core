package log

import (
	"demo/tools/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"math"
	"os"
	"path/filepath"
	"time"
)

var (
	_minTimeInt64 = time.Unix(0, math.MinInt64)
	_maxTimeInt64 = time.Unix(0, math.MaxInt64)

	exportLogger      *zap.Logger
	exportSugarLogger *zap.SugaredLogger
)

const (
	flagLevel = "log.level"
)

const DefaultLogPath = "/var/log/test" // 默认输出日志文件路径

// getWriter 自定义Writer,分割日志
func getWriter(conf *config.Log) zapcore.WriteSyncer {
	if len(conf.OutputPaths) < 1 {
		conf.OutputPaths = []string{DefaultLogPath}
	}
	// 判断日志路径是否存在，如果不存在就创建
	var curPath = ""
	for _, path := range conf.OutputPaths {
		if exist := IsExist(path); !exist {
			if err := os.MkdirAll(path, os.ModePerm); err != nil {
				panic(any("log path is nil"))
				return nil
			} else {
				curPath = path
				break
			}
		}
	}
	// 日志文件 与 日志切割 配置
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filepath.Join(curPath, conf.Filename), // 日志文件路径
		MaxSize:    conf.Maxsize,                          // 单个日志文件最大多少 mb
		MaxBackups: 1,                                     // 日志备份数量
		MaxAge:     conf.MaxAge,                           // 日志最长保留时间
		Compress:   conf.Compress,                         // 是否压缩日志
	}
	if config.IsDev() {
		// 日志同时输出到控制台和日志文件中
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(lumberJackLogger), zapcore.AddSync(os.Stdout))
	}
	// 日志只输出到日志文件
	return zapcore.AddSync(lumberJackLogger)
}

// IsExist 判断文件或者目录是否存在
func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// getProdEncoder 自定义日志编码器
func getProdEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

//func getDevEncoder() zapcore.Encoder {
//	encoderConfig := zap.NewDevelopmentEncoderConfig()
//	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
//	return zapcore.NewConsoleEncoder(encoderConfig)
//}

func getDevEncoder() zapcore.Encoder {
	encoderConfig := newDevelopmentEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func newDevelopmentEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		// Keys can be anything except the empty string.
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     CustomTimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

// CustomTimeEncoder 时间格式修改
func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02T15:04:05"))
}

func getLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		panic(any("log level error"))
	}
}
