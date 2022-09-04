package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var zapLog *zap.Logger

func InitializeLogger() {
	var err error
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	zapLog, err = config.Build()
	if err != nil {
		panic(err)
	}
}

func Info(message string, fields ...zap.Field) {
	zapLog.Info(message, fields...)
	zapLog.Sync()
}

func Debug(message string, fields ...zap.Field) {
	zapLog.Debug(message, fields...)
	zapLog.Sync()
}

func Error(message string, fields ...zap.Field) {
	zapLog.Error(message, fields...)
	zapLog.Sync()
}

func Fatal(message string, fields ...zap.Field) {
	zapLog.Fatal(message, fields...)
	zapLog.Sync()
}

func Warn(message string, fields ...zap.Field) {
	zapLog.Warn(message, fields...)
	zapLog.Sync()
}

func Panic(message string, fields ...zap.Field) {
	zapLog.Panic(message, fields...)
	zapLog.Sync()
}
