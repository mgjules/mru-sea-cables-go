package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New returns a new logger
func New(dev, debug bool) *zap.SugaredLogger {
	atom := zap.NewAtomicLevel()

	switch {
	case debug:
		atom.SetLevel(zapcore.DebugLevel)
	case dev:
		atom.SetLevel(zapcore.InfoLevel)
	default:
		atom.SetLevel(zapcore.WarnLevel)
	}

	logger := zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
		zapcore.Lock(os.Stdout),
		atom,
	))

	return logger.Sugar()
}
