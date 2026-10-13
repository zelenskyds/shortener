package config

import "fmt"

var BaseHost = "localhost"
var BasePort = "8080"

var BaseURL = fmt.Sprintf("%s:%s", BaseHost, BasePort)
var BaseGetURL = BaseURL
