package log

import (
	"demo/tools/config"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

var logger *zap.Logger

func init() {
	var err error
	logger, err = New(config.Conf.Log)
	if err != nil {
		_ = fmt.Errorf("err init new logger %v", err)
		return
	}
	loggerF = logger.Sugar()
	logger.Info("Logger initialization successful")
}
func New(cfg *config.Log) (*zap.Logger, error) {
	var (
		level    zapcore.Level
		writer   zapcore.WriteSyncer
		encoder  zapcore.Encoder
		core     zapcore.Core
		fileCore zapcore.Core
	)

	// 解析日志级别
	if err := level.UnmarshalText([]byte(cfg.Level)); err != nil {
		return nil, err
	}

	// 设置日志编码器
	switch cfg.Encoding {
	case "json":
		encoder = zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	case "console":
		encoder = zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig())
	case "dev":
		encoder = getDevEncoder()
	case "prod":
		encoder = getProdEncoder()
	default:
		return nil, errors.New("invalid encoding")
	}
	// 设置日志输出
	if cfg.Stdout {
		writer = zapcore.AddSync(os.Stdout)
	}
	if cfg.Filename != "" {
		fileWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:   cfg.Filename,
			MaxSize:    cfg.Maxsize,
			MaxAge:     cfg.MaxAge,
			MaxBackups: cfg.FileMaxBackups,
			LocalTime:  true,
			Compress:   cfg.Compress,
		})
		fileCore = zapcore.NewCore(
			encoder,
			fileWriter,
			level,
		)
		if writer != nil {
			writer = zapcore.NewMultiWriteSyncer(writer, fileWriter)
		} else {
			writer = fileWriter
		}
	}
	// 组合日志核心
	core = zapcore.NewCore(encoder, writer, level)

	// 添加 Caller 和 StackTrace
	newLogger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.DPanicLevel))
	if fileCore != nil {
		fileLogger := zap.New(fileCore, zap.AddCaller(), zap.AddStacktrace(zapcore.DPanicLevel))
		defer fileLogger.Sync()
		fileLogger.Info("fileLogger initialization successful")
	}
	return newLogger, nil
}
