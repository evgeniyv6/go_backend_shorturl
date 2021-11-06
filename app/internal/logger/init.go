package logger

import (
	"os"

	"go.elastic.co/ecszap"
	"go.uber.org/zap"
)

type (
	ZapLogger struct {
		logger *zap.SugaredLogger
	}

	ZapWrapper interface {
		Error(msg string)
		Info(args ...interface{})
		Infof(msg string, args ...interface{})
		Infow(msg string, keysAndValues ...interface{})
		Errorw(msg string, keysAndValues ...interface{})
		Panicw(msg string, keysAndValues ...interface{})
	}
)

func NewZapWrapper() ZapWrapper {
	encoderConfig := ecszap.NewDefaultEncoderConfig()
	core := ecszap.NewCore(encoderConfig, os.Stdout, zap.DebugLevel)
	logger := zap.New(core, zap.AddCaller())
	logger = logger.With(zap.String("app", "link cutter"))
	sugar := logger.Sugar()

	defer func() { _ = logger.Sync() }()

	return &ZapLogger{logger: sugar}
}

func (z *ZapLogger) Error(msg string) {
	z.logger.Error(msg)
}
func (z *ZapLogger) Infof(msg string, args ...interface{}) {
	z.logger.Infof(msg, args...)
}

func (z *ZapLogger) Panicw(msg string, keysAndValues ...interface{}) {
	z.logger.Panicw(msg, keysAndValues)
}

func (z *ZapLogger) Info(args ...interface{}) {
	z.logger.Info(args)
}

func (z *ZapLogger) Errorw(msg string, keysAndValues ...interface{}) {
	z.logger.Errorw(msg, keysAndValues)
}

func (z *ZapLogger) Infow(msg string, keysAndValues ...interface{}) {
	z.logger.Infow(msg, keysAndValues)
}
