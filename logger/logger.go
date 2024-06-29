package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger    *zap.SugaredLogger
	zapLogger *zap.Logger
)

func init() {
	zapLogger, _ = zap.NewProduction()
	logger = zapLogger.Sugar()
}

func GetLogger() *zap.SugaredLogger {
	return logger
}

func InitLogger(env, logLevel string) {
	var level zap.AtomicLevel
	err := level.UnmarshalText([]byte(logLevel))
	if err != nil {
		logger.Fatalf("invalid log level: %v", err)
	}
	zapConfig := &zap.Config{
		Level:       level,
		Development: true,
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "message",
			LevelKey:       "level",
			TimeKey:        "time",
			CallerKey:      "caller",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		InitialFields:    nil,
	}

	if env == "prod" {
		zapConfig.Development = false
	}

	zapLogger, err = zapConfig.Build()
	if err != nil {
		logger.Fatal("zap config could not be built: %v", err)
	}
	*logger = *zapLogger.Sugar()
}
