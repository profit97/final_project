package main

import (
	"log"
	"os"
	"strconv"

	"il.karabach/pkg/db"
	"il.karabach/pkg/server"
)

func main() {
	// Определяем порт веб-сервера 7540
	port := 7540
	if portStr := os.Getenv("TODO_PORT"); portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			port = p
		} else {
			log.Printf("Некорректное значение TODO_PORT %q, используется порт %d", portStr, port)
		}
	}

	// Инициализация SQLite
	dbPath := "scheduler.db" 
	database, err := db.InitDB(dbPath)
	if err != nil {
		log.Fatalf("Ошибка инициализации БД: %v", err)
	}
	defer database.Close()

	// Запуск сервера (статических файлов)
	if err := server.Start(port, "./web", database); err != nil {
		log.Fatal(err)
	}
}