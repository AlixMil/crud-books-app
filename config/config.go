package config

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	ServerPort          string
	GoFileServiceApiKey string
	GoFileFolderToken   string
	DatabaseName        string
	DatabasePort        string
	DatabaseHost        string
	DatabaseLogin       string
	DatabasePwd         string
	JwtSecret           string
	AccessTokenTTL      time.Duration
	RefreshTokenTTL     time.Duration
}

func New() (*Config, error) {
	AccessTokenTTL, err := time.ParseDuration(os.Getenv("ACCESS_TOKEN_TTL"))
	if err != nil {
		return nil, fmt.Errorf("parse access token duration: %w", err)
	}

	RefreshTokenTTL, err := time.ParseDuration(os.Getenv("REFRESH_TOKEN_TTL"))
	if err != nil {
		return nil, fmt.Errorf("parse refresh token duration: %w", err)
	}

	cfg := Config{
		ServerPort:          os.Getenv("SERVER_PORT"),
		GoFileServiceApiKey: os.Getenv("GOFILE_SERVICE_API_KEY"),
		GoFileFolderToken:   os.Getenv("GOFILE_FOLDER_TOKEN"),
		DatabaseName:        os.Getenv("DB_NAME"),
		DatabasePort:        os.Getenv("DB_PORT"),
		DatabaseHost:        os.Getenv("DB_HOST"),
		DatabaseLogin:       os.Getenv("DB_LOGIN"),
		DatabasePwd:         os.Getenv("DB_PWD"),
		JwtSecret:           os.Getenv("JWT_SECRET"),
		AccessTokenTTL:      AccessTokenTTL,
		RefreshTokenTTL:     RefreshTokenTTL,
	}

	return &cfg, nil
}

func (c *Config) Validate() error {
	if c.ServerPort == "" {
		return fmt.Errorf("serverport env is empty")
	}
	if c.GoFileServiceApiKey == "" {
		return fmt.Errorf("goFileServiceApiKey env is empty")
	}
	if c.GoFileFolderToken == "" {
		return fmt.Errorf("goFileFolderToken env is empty")
	}
	if c.DatabaseName == "" {
		return fmt.Errorf("databaseName env is empty")
	}
	if c.DatabasePort == "" {
		return fmt.Errorf("databasePort env is empty")
	}
	if c.DatabaseHost == "" {
		return fmt.Errorf("databaseHost env is empty")
	}
	if c.DatabaseLogin == "" {
		return fmt.Errorf("databaseLogin env is empty")
	}
	if c.JwtSecret == "" {
		return fmt.Errorf("databasePwd env is empty")
	}
	if c.DatabaseHost == "" {
		return fmt.Errorf("jwtSecret env is empty")
	}
	if c.AccessTokenTTL == 0 {
		return fmt.Errorf("accessTokenTTL env is empty or null")
	}
	if c.RefreshTokenTTL == 0 {
		return fmt.Errorf("refreshTokenTTL env is empty or null")
	}
	return nil
}
