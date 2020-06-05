package main

import (
	"dynamic-protobuf-json-service/config"
	"dynamic-protobuf-json-service/engine"
	"dynamic-protobuf-json-service/log"
	"fmt"
)

func main() {

	var err error
	cfg, err := config.GetAppConfig()
	if err != nil {
		er(fmt.Errorf("Could not parse application configuration, reason: %v", err.Error()))
		return
	}

	logger, err := log.InitLogger(cfg)
	if err != nil {
		er(fmt.Errorf("Could not initialize logger, reason: %v", err.Error()))
		return
	}
	defer logger.Sync()

	if err = engine.Run(cfg); err != nil {
		er(fmt.Errorf("Could not run engine, reason: %v", err.Error()))
		return
	}
}

func er(err error) {
	fmt.Println(err.Error())
}
