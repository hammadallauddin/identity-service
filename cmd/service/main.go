package main

import (
	"log/slog"
	"os"
)

func main() {
	handlerOpts := &slog.HandlerOptions{Level: slog.LevelDebug}
	logger := slog.New(slog.NewTextHandler(os.Stderr, handlerOpts))
	logger.Info("Identity Service starting up")
	logger = logger.With("service", "identity-service")

}
