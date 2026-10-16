package main

import (
	"net/http"

	"github.com/zelenskyds/shortener/internal/config"
	"github.com/zelenskyds/shortener/internal/handler"
	"github.com/zelenskyds/shortener/internal/logger"
	"github.com/zelenskyds/shortener/internal/repository"
	"go.uber.org/zap"
)

func main() {
	ParseFlag()

	err := run()

	if err != nil {
		panic(err)
	}
}

func run() error {
	db := &repository.Storage{}

	if err := logger.Initialize(config.LogLevel); err != nil {
		return err
	}

	logger.Log.Info("Running server", zap.String("address", config.ServerAddress))

	return http.ListenAndServe(config.ServerAddress, logger.HTTPHandlerLogger(handler.NewRouter(config.BaseURL, db)))
}
