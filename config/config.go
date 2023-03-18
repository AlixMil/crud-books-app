package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerHost           string
	ServerPort           string
	StorageServiceApiKey string
	DatabaseName         string
	DatabasePort         string
	DatabaseHost         string
	DatabaseLogin        string
	DatabasePwd          string
}

func New() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return &Config{}, err
	}

	return &Config{
		ServerHost:           os.Getenv("SERVER_HOST"),
		ServerPort:           os.Getenv("SERVER_PORT"),
		StorageServiceApiKey: os.Getenv("STORAGE_SERVICE_API_KEY"),
		DatabaseName:         os.Getenv("DB_NAME"),
		DatabasePort:         os.Getenv("DB_PORT"),
		DatabaseHost:         os.Getenv("DB_HOST"),
		DatabaseLogin:        os.Getenv("DB_LOGIN"),
		DatabasePwd:          os.Getenv("DB_PWD"),
	}, nil
}
