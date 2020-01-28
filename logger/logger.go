package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New returns a new logger
func New(dev, verbose bool) *zap.SugaredLogger {
	var encoderCfg zapcore.EncoderConfig
	if dev {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderCfg = zap.NewProductionEncoderConfig()
	}

	atom := zap.NewAtomicLevel()

	switch {
	case dev && verbose:
		atom.SetLevel(zapcore.DebugLevel)
	case dev && !verbose:
		atom.SetLevel(zapcore.InfoLevel)
	case !dev && verbose:
		atom.SetLevel(zapcore.WarnLevel)
	case dev && !verbose:
		atom.SetLevel(zapcore.ErrorLevel)
	}

	logger := zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atom,
	))

	return logger.Sugar()
}
