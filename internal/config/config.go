package config

import (
	"os"
)

type Config struct {
	TgToken    string //telegram auth token
	BitopToken string //bitop auth token
	SchPath    string //schedule path
	SgPath     string //group search path
	DbAuth     DBAuthData
}

type DBAuthData struct {
	DBAdminUsername string
	DBAdminPassword string
	DBName          string
	DBHost          string
	DBPort          string
	URI             string
}

func New() *Config {
	return &Config{
		TgToken:    os.Getenv("TG_TOKEN"),
		BitopToken: os.Getenv("BITOP_TOKEN"),
		SchPath:    os.Getenv("SCHEDULE_PATH"),
		SgPath:     os.Getenv("SG_PATH"),
		DbAuth: DBAuthData{
			DBAdminUsername: os.Getenv("DB_ADMIN_USERNAME"),
			DBAdminPassword: os.Getenv("DB_ADMIN_PASSWORD"),
			DBHost:          os.Getenv("DB_HOST"),
			DBPort:          os.Getenv("DB_PORT"),
			DBName:          os.Getenv("DB_NAME"),
			URI:             os.Getenv("POSTGRES_URI"),
		},
	}
}
