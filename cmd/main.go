package main

import (
	"crud-books/gofile"
	"crud-books/handlers"
	"crud-books/mongodb"
	"crud-books/server"
	"crud-books/storage"
	"log"
)

func main() {
	db, err := mongodb.New("user", "password", "dsn")
	if err != nil {
		log.Fatalf("Database init failed: %v", err)
	}

	service := gofile.New("API_KEY")
	strg, err := storage.New(service)
	if err != nil {
		log.Fatalf("storage error: %v", err)
	}

	hndlrs, err := handlers.New(db, strg)
	if err != nil {
		log.Fatalf("handlers init: %v", err)
	}

	srv := server.New("4001")
	srv.InitHandlers(hndlrs)
	if err := srv.Start(); err != nil {
		log.Fatalf("server is not started. Error: %v", err)
	}
}
