package main

import (
	"github.com/zelenskyds/shortener/internal/config"
	"github.com/zelenskyds/shortener/internal/handler"
)

func main() {
	ParseFlag()
	handler.RunApp(config.ServerAddress, config.BaseURL)
}
