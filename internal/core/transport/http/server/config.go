package core_http_server

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Addr            string
	ShutdownTimeout time.Duration
}

func LoadServerConfig() Config {
	val := os.Getenv("HTTP_TIMEOUT")
	i, err := strconv.Atoi(val)
	if err != nil {
		panic("invalid HTTP_TIMEOUT")
	}
	timeout := time.Duration(i) * time.Second
	return Config{
		Addr:            os.Getenv("HTTP_ADDRESS"),
		ShutdownTimeout: timeout,
	}
}
