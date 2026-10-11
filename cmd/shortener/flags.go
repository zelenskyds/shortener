package main

import (
	"flag"
	"strings"

	"github.com/zelenskyds/shortener/internal/config"
)

func ParseFlag() {
	flag.StringVar(&config.BaseURL, "a", config.BaseURL, "установить адрес")
	flag.StringVar(&config.BaseGetURL, "b", config.BaseGetURL, "установить адрес запроса урла")
	flag.Parse()

	config.BaseURL = strings.TrimPrefix(config.BaseURL, "http://")
	config.BaseURL = strings.TrimPrefix(config.BaseURL, "https://")
	config.BaseGetURL = strings.TrimPrefix(config.BaseGetURL, "http://")
	config.BaseGetURL = strings.TrimPrefix(config.BaseGetURL, "https://")
}
