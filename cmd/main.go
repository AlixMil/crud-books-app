package main

import (
	"crud-books/auth"
	"crud-books/config"
	"crud-books/handlers"
	"crud-books/mongodb"
	"crud-books/server"
	"crud-books/services"
	"crud-books/storage"
	"log"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("config initializing: %v", err)
	}
	err = cfg.Validate()
	if err != nil {
		log.Fatalf("config validation: %v", err)
	}

	db := mongodb.New()
	err = db.Connect(cfg)
	if err != nil {
		log.Fatalf("database init: %v", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("database ping: %v", err)
	}

	storage := storage.New(cfg)

	jwtEngine := auth.NewJwtEngine(cfg)
	hashEngine := auth.NewHashEngine()

	services := services.New(db, jwtEngine, storage, hashEngine)

	handlers := handlers.New(services)
	if err != nil {
		log.Fatalf("handlers init: %v", err)
	}

	srv := server.New(
		cfg.ServerPort,
		cfg.JwtSecret,
		jwtEngine.GetMiddleware(),
	)
	srv.InitMiddlewares()
	srv.UseRouters(handlers)

	err = srv.Start()
	if err != nil {
		log.Fatalf("server isn't started: %v", err)
	}
}
