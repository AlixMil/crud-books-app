package main

import (
	"crud-books/config"
	"crud-books/handlers"
	"crud-books/mongodb"
	"crud-books/server"
	storageService "crud-books/storageService"
	"log"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	db, err := mongodb.New(cfg.DatabaseLogin, cfg.DatabasePwd, cfg.DatabaseName, cfg.DatabaseHost, cfg.DatabasePort)
	if err != nil {
		log.Fatalf("Database init failed: %v", err)
	}

	storage := storageService.New(cfg.StorageServiceApiKey)

	// strg, err := storage.New(service)
	// if err != nil {
	// 	log.Fatalf("storage error: %v", err)
	// }

	handlers, err := handlers.New(db, storage)
	if err != nil {
		log.Fatalf("handlers init: %v", err)
	}

	srv := server.New("4001")
	srv.InitHandlers(handlers)
	if err := srv.Start(); err != nil {
		log.Fatalf("server is not started. Error: %v", err)
	}
}
