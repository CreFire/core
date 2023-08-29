package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"testing"
)

func BenchmarkLog(b *testing.B) {
	b.ResetTimer()
	_, err := NewDefault()
	if err != nil {
		return
	}
	for i := 0; i < b.N; i++ {
		InfoF("A walrus appears", "walrus",
			1,
			1.01)
	}
}

func BenchmarkZap(b *testing.B) {

	loges := zap.New(
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
			zapcore.AddSync(io.Discard),
			zap.InfoLevel),
	)
	l := loges.Sugar()
	//logs.Core().With([]zap.Field{String("k", "v")})
	defer loges.Sync()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		l.Info("A walrus appears", "walrus",
			1,
			1.01)
	}
}

//
//func BenchmarkLogrus(b *testing.B) {
//	logger := logrus.New()
//	logger.SetLevel(logrus.DebugLevel)
//	logger.SetOutput(io.Discard)
//
//	b.ResetTimer()
//
//	for i := 0; i < b.N; i++ {
//		logger.WithFields(logrus.Fields{
//			"animal": "walrus",
//			"number": 1,
//			"size":   10.1,
//		}).Debug("A walrus appears")
//	}
//}
