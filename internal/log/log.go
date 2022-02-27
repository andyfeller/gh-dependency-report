package log

import (
	"go.uber.org/zap"
)

func NewLogger(debug bool) (*zap.Logger, error) {

	level := zap.InfoLevel

	if debug {
		level = zap.DebugLevel
	}

	loggerConfig := zap.Config{
		Level:            zap.NewAtomicLevelAt(level),
		Encoding:         "console",
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	return loggerConfig.Build()
}
