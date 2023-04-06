package main

import (
	"crud-books/config"
	"crud-books/handlers"
	"crud-books/mongodb"
	"crud-books/pkg/hasher"
	jwt_package "crud-books/pkg/jwt"
	"crud-books/server"
	"crud-books/services"
	"crud-books/storageService/gofile"
	"log"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	db, err := mongodb.NewClient(*cfg)
	if err != nil {
		log.Fatalf("Database init failed: %v", err)
	}

	storage := gofile.New(cfg.GoFileServiceApiKey, cfg.GoFileFolderToken)

	tokener := jwt_package.New(*cfg)

	hasher := hasher.New()

	services := services.New(db, *tokener, *storage, *hasher)

	handlers, err := handlers.New(db, storage, cfg.JWTSecret, services, tokener)
	if err != nil {
		log.Fatalf("handlers init: %v", err)
	}

	srv := server.New("4001", cfg.JWTSecret)
	srv.InitHandlers(handlers)
	if err := srv.Start(); err != nil {
		log.Fatalf("server is not started. Error: %v", err)
	}
}
