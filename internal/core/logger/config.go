package logger

import (
	"os"
)

type Config struct {
	Level  string
	Folder string
}

func LoadConfig() Config {
	cfg := Config{
		Level:  os.Getenv("LOGGER_LEVEL"),
		Folder: os.Getenv("LOGGER_FOLDER"),
	}

	if cfg.Level == "" {
		panic("Logger Level not set")
	}
	if cfg.Folder == "" {
		panic("Logger Folder not set")
	}
	return cfg
}
