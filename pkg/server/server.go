package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"il.karabach/pkg/api" 
)

// Start запускает HTTP-сервер с файловым сервером для статики и API.
func Start(port int, webDir string, db *sql.DB) error {
	// Инициализация API 
	api.Init()

	// Файловый сервер для статических файлов
	fs := http.FileServer(http.Dir(webDir))
	http.Handle("/", fs)

	addr := fmt.Sprintf(":%d", port)
	log.Printf("Сервер запущен на http://localhost%s", addr)
	return http.ListenAndServe(addr, nil)
}