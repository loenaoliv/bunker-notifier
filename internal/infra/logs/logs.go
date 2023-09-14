package logs

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(logLevel zapcore.Level) *zap.Logger {
	logger, err := configureLogger(logLevel)
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
	return zap.L()
}

func configureLogger(logLevel zapcore.Level) (*zap.Logger, error) {
	var zapConfig zap.Config

	encodingConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		MessageKey:     "msg",
		CallerKey:      "caller",
		StacktraceKey:  "stackTrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	}

	zapConfig = zap.Config{
		Level:            zap.NewAtomicLevelAt(logLevel),
		Encoding:         "json",
		EncoderConfig:    encodingConfig,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := zapConfig.Build()
	if err != nil {
		return nil, err
	}

	return logger, nil
}
