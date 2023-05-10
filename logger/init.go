package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var NomiosLog *Logger

func init() {
	InitLogger()
}

// InitLogger -- init logger
func InitLogger() {
	// init production encoder config
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	// init production config
	cfg := zap.NewProductionConfig()
	cfg.Sampling = nil
	cfg.EncoderConfig = encoderCfg
	cfg.OutputPaths = []string{"stdout"}
	cfg.ErrorOutputPaths = []string{"stdout"}
	// build logger
	logger, _ := cfg.Build()
	logger = logger.WithOptions(
		zap.AddCallerSkip(1),
	)

	NomiosLog = &Logger{
		sl: logger.Sugar(),
	}
}
