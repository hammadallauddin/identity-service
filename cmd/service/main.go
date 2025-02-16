package main

import (
	"log"

	"github.com/hammadallauddin/identity-service/pkg/config"
	"github.com/hammadallauddin/identity-service/pkg/logs"
)

func main() {
	err := config.Initialize()
	if err != nil {
		log.Fatalf("could not initialize config: %v", err)
	}

	logger, err := logs.Initialize()
	if err != nil {
		log.Fatalf("could not initialize logger: %v", err)
	}

	logger.Debug("starting service ... ")
}
