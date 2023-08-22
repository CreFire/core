package log

import (
	"go.uber.org/zap"
)

var loggerF *zap.SugaredLogger

func DebugF(msg string, fields ...any) {
	loggerF.Debug(msg, fields)
}

func InfoF(msg string, fields ...any) {
	loggerF.Info(msg, fields)
}

func WarnF(msg string, fields ...any) {
	loggerF.Warn(msg, fields)
}

func ErrorF(msg string, fields ...any) {
	loggerF.Error(msg, fields)
}

func DPanicF(msg string, fields ...any) {
	loggerF.DPanic(msg, fields)
}

func PanicF(msg string, fields ...any) {
	loggerF.Panic(msg, fields)
}

func FatalF(msg string, fields ...any) {
	loggerF.Fatal(msg, fields)
}
