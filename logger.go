package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"pirate-wars/cmd/common"
)

func createLogger() *zap.SugaredLogger {
	// truncate file
	configFile, err := os.OpenFile(common.LogFile, os.O_TRUNC|os.O_CREATE, 0664)
	if err != nil {
		panic(err)
	}
	if err = configFile.Close(); err != nil {
		panic(err)
	}
	// create logger
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{common.LogFile}
	cfg.Level = zap.NewAtomicLevelAt(BASE_LOG_LEVEL)
	cfg.Development = DEV_MODE
	cfg.DisableCaller = false
	cfg.DisableStacktrace = false
	cfg.Encoding = "console"
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncoderConfig = encoderConfig
	logger := zap.Must(cfg.Build())
	defer logger.Sync()
	return logger.Sugar()
}
