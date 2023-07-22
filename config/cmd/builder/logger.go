package builder

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
)

// NewLogger defines all configurations to instantiate a log.
func NewLogger() *zap.SugaredLogger {
	config := zap.NewProductionConfig()
	config.Encoding = "console"
	config.DisableCaller = true
	config.DisableStacktrace = true
	config.EncoderConfig.TimeKey = "time"
	config.EncoderConfig.LevelKey = "level"
	config.EncoderConfig.MessageKey = "msg"
	config.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.ConsoleSeparator = "  "

	logger, err := config.Build()
	if err != nil {
		log.Fatalf("Error creating logger: %v", err)
	}
	return logger.Sugar()
}

// Sync defines the synchronization between generated logs.
func Sync(log *zap.SugaredLogger) {
	_ = log.Sync()
}
