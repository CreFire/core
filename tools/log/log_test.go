package log

import (
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"testing"
)

func BenchmarkLog(b *testing.B) {
	logs, _ := NewDefault()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logs.Info("A")
	}
}

func BenchmarkZap(b *testing.B) {

	logs := zap.New(
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
			zapcore.AddSync(io.Discard),
			zap.InfoLevel),
	)
	logs.Core().With([]zap.Field{String("k", "v")})
	defer logs.Sync()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logs.Info("A walrus appears")
	}
}

func BenchmarkLogrus(b *testing.B) {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.SetOutput(io.Discard)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.WithFields(logrus.Fields{
			"animal": "walrus",
			"number": 1,
			"size":   10.1,
		}).Debug("A walrus appears")
	}
}
