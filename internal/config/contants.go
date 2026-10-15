package config

import "fmt"

var BaseHost = "localhost"
var BasePort = "8080"

var ServerAddress = fmt.Sprintf("%s:%s", BaseHost, BasePort)
var BaseURL = fmt.Sprintf("http://%s", ServerAddress)
