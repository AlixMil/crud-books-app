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

	db := mongodb.New()
	err = db.Connect(*cfg)
	if err != nil {
		log.Fatalf("database init failed: %v", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("database ping execute error")
	}

	storage := gofile.New(cfg.GoFileServiceApiKey, cfg.GoFileFolderToken)

	tokener := jwt_package.New(*cfg)

	hasher := hasher.New()

	services := services.New(db, *tokener, *storage, *hasher)

	handlers := handlers.New(services)
	if err != nil {
		log.Fatalf("handlers init: %v", err)
	}

	srv := server.New("4001", cfg.JWTSecret)
	srv.InitMiddlewares()
	srv.UseRouters(handlers)
	if err := srv.Start(); err != nil {
		log.Fatalf("server is not started. Error: %v", err)
	}
}
