package log

import "go.uber.org/zap"

func Debug(msg string, fields ...zap.Field) {
	defLog.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	defLog.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	defLog.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	defLog.Error(msg, fields...)
}

func DPanic(msg string, fields ...zap.Field) {
	defLog.DPanic(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	defLog.Panic(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	defLog.Fatal(msg, fields...)
}
