package log

import (
	"fmt"
	"go.uber.org/zap"
)

var loggerF *zap.SugaredLogger

func DebugF(msg string, fields ...any) {
	str := fmt.Sprint(fields...)
	loggerF.Debug(str)
}

func InfoF(fields ...any) {
	str := fmt.Sprint(fields...)
	loggerF.Info(str)
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
