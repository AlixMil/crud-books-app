package config

import (
	"errors"
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

	cfg := Config{
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
	}

	if cfg.ServerHost == "" {
		return nil, errors.New("serverhost env is empty")
	}
	if cfg.ServerPort == "" {
		return nil, errors.New("serverport env is empty")
	}
	if cfg.DdownloadServiceApiKey == "" {
		return nil, errors.New("ddownloadServiceApiKey env is empty")
	}
	if cfg.GoFileServiceApiKey == "" {
		return nil, errors.New("goFileServiceApiKey env is empty")
	}
	if cfg.GoFileFolderToken == "" {
		return nil, errors.New("goFileFolderToken env is empty")
	}
	if cfg.DatabaseName == "" {
		return nil, errors.New("databaseName env is empty")
	}
	if cfg.DatabasePort == "" {
		return nil, errors.New("databasePort env is empty")
	}
	if cfg.DatabaseHost == "" {
		return nil, errors.New("databaseHost env is empty")
	}
	if cfg.DatabaseLogin == "" {
		return nil, errors.New("databaseLogin env is empty")
	}
	if cfg.JWTSecret == "" {
		return nil, errors.New("databasePwd env is empty")
	}
	if cfg.DatabaseHost == "" {
		return nil, errors.New("jwtSecret env is empty")
	}
	if cfg.JWTTokenTTL == 0 {
		return nil, errors.New("jwtTokenTTL env is empty or null")
	}

	return &cfg, nil
}
