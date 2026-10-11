package config

import "fmt"

var BASE_HOST = "localhost"
var BASE_PORT = "8080"

var BASE_URL = fmt.Sprintf("%s:%s", BASE_HOST, BASE_PORT)
