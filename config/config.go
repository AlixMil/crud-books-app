package config

import (
	"fmt"
	"os"
	"strconv"

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
	JWTTokenTTL            int
}

func New() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return &Config{}, err
	}

	JWTTokenTTL, err := strconv.Atoi(os.Getenv("JWT_TOKEN_TTL"))
	if err != nil {
		return nil, fmt.Errorf("failed convert to in of JWT_TOKEN_TTL, error: %w", err)
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
		JWTTokenTTL:            JWTTokenTTL,
	}, nil
}
