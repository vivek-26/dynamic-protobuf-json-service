package log

import (
	"dynamic-protobuf-json-service/config"

	"go.uber.org/zap"
)

// InitLogger ...
func InitLogger(cfg *config.Config) (*zap.Logger, error) {

	var logger *zap.Logger
	var err error

	switch cfg.Env {
	case "development":
		logger, err = zap.NewDevelopment()
	case "production":
		logger, err = zap.NewProduction()
	default:
		logger = zap.NewExample()
	}

	if err != nil {
		return nil, err
	}

	zap.ReplaceGlobals(logger)
	return logger, nil
}
