package config

import (
	"os"
)

type Config struct {
	TgToken    string //telegram auth token
	BitopToken string //bitop auth token
	SchPath    string //schedule path
	SgPath     string //group search path
}

func New() *Config {
	return &Config{
		TgToken:    os.Getenv("TG_TOKEN"),
		BitopToken: os.Getenv("BITOP_TOKEN"),
		SchPath:    os.Getenv("SCHEDULE_PATH"),
		SgPath:     os.Getenv("SG_PATH"),
	}
}
