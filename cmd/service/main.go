package main

import (
	"log/slog"

	"github.com/hammadallauddin/identity-service/pkg/config"
)

func main() {
	err := config.Initialize()
	if err != nil {
		slog.Error("unable to initialze config", err)
	}
	slog.Info("starting service")
}
