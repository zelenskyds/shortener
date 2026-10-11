package main

import (
	"github.com/zelenskyds/shortener/internal/config"
	"github.com/zelenskyds/shortener/internal/handler"
)

func main() {
	handler.RunApp(config.BASE_URL)
}
