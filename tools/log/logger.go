package log

import (
	"core/tools/config"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
)

var defLog *zap.Logger

func init() {
	var err error
	defLog, err = New(config.Conf.Log)
	if err != nil {
		_ = fmt.Errorf("err init new defLog %v", err)
		return
	}
	loggerF = defLog.Sugar()
	loggerF.WithOptions(AddCallerSkip(1))
	Info("Logger initialization successful")
}
func NewDefault() (*zap.Logger, error) {
	var (
		writer  zapcore.WriteSyncer
		encoder zapcore.Encoder
		core    zapcore.Core
	)

	// 设置日志编码器
	encoder = zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	// 设置日志输出
	writer = zapcore.AddSync(io.Discard)
	// 组合日志核心
	core = zapcore.NewCore(encoder, writer, zapcore.InfoLevel)

	defLog = zap.New(core) // zap.AddCaller(), zap.AddStacktrace(zapcore.DPanicLevel)
	defLog = defLog.WithOptions(AddCallerSkip(1))
	loggerF = defLog.Sugar()
	return defLog, nil
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
	newLogger = newLogger.WithOptions(AddCallerSkip(1))
	if fileCore != nil {
		fileLogger := zap.New(fileCore, zap.AddCaller(), zap.AddStacktrace(zapcore.DPanicLevel))
		defer fileLogger.Sync()
		fileLogger.Info("fileLogger initialization successful")
	}
	return newLogger, nil
}
