package config

import (
	"encoding/json"
	"go.uber.org/zap"
	"time"
)

func Logger() *zap.Logger {
	Init()
	vipr := GetConfig()
	rawJSON := []byte(`{
   "level": "` + vipr.GetString("LOG_LEVEL") + `",
   "encoding": "json",
   "outputPaths": ["stdout"],
   "errorOutputPaths": ["stderr"],
   "encoderConfig": {
     "messageKey": "message",
     "levelKey": "level",
     "levelEncoder": "lowercase"
   }
 }`)
	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	sugar := logger.Sugar()
	sugar.Infow("Failed to fetch URL.",
		// Structured context as loosely typed key-value pairs.
		"intent", "get token from mpesa",
		"attempt", 1,
		"backoff", time.Second,
	)
	var execution = "get/mpesa/token"
	sugar.Infof("Failed to execute: %s", execution)
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)

	return logger
}
