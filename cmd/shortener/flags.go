package main

import (
	"flag"
	"os"
	"strings"

	"github.com/zelenskyds/shortener/internal/config"
)

func ParseFlag() {
	flag.StringVar(&config.ServerAddress, "a", config.ServerAddress, "установить адрес")
	flag.StringVar(&config.BaseURL, "b", config.BaseURL, "установить адрес запроса урла")
	flag.Parse()

	if serverAddress := os.Getenv("SERVER_ADDRESS"); serverAddress != "" {
		config.ServerAddress = serverAddress
	}

	if baseURL := os.Getenv("BASE_URL"); baseURL != "" {
		config.BaseURL = baseURL
	}

	config.ServerAddress = trimScheme(config.ServerAddress)
	config.BaseURL = withScheme(config.BaseURL)
}

func withScheme(address string) string {
	if strings.Contains(address, "://") {
		return address
	}

	return "http://" + address
}

func trimScheme(address string) string {
	address = strings.TrimPrefix(address, "http://")
	address = strings.TrimPrefix(address, "https://")

	return address
}
