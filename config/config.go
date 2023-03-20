package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerHost             string
	ServerPort             string
	DdownloadServiceApiKey string
	GoFileServiceApiKey    string
	GoFileFolderToken      string
	DatabaseName           string
	DatabasePort           string
	DatabaseHost           string
	DatabaseLogin          string
	DatabasePwd            string
	JWTSecret              string
}

func New() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return &Config{}, err
	}

	return &Config{
		ServerHost:             os.Getenv("SERVER_HOST"),
		ServerPort:             os.Getenv("SERVER_PORT"),
		DdownloadServiceApiKey: os.Getenv("DDOWNLOAD_SERVICE_API_KEY"),
		GoFileServiceApiKey:    os.Getenv("GOFILE_SERVICE_API_KEY"),
		GoFileFolderToken:      os.Getenv("GOFILE_FOLDER_TOKEN"),
		DatabaseName:           os.Getenv("DB_NAME"),
		DatabasePort:           os.Getenv("DB_PORT"),
		DatabaseHost:           os.Getenv("DB_HOST"),
		DatabaseLogin:          os.Getenv("DB_LOGIN"),
		DatabasePwd:            os.Getenv("DB_PWD"),
		JWTSecret:              os.Getenv("JWT_SECRET"),
	}, nil
}
