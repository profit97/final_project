package db

import (
	"database/sql"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB(dbPath string) (*sql.DB, error) {
	needCreate := false
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		needCreate = true
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	if needCreate {
		log.Println("Файл базы данных не найден. Создаём новую БД и таблицу scheduler.")
		if err := createTable(db); err != nil {
			db.Close()
			return nil, err
		}
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	DB = db
	log.Println("База данных SQLite успешно инициализирована")
	return db, nil
}

func createTable(db *sql.DB) error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date CHAR(8) NOT NULL DEFAULT "",
		title TEXT NOT NULL,
		comment VARCHAR(64),
		repeat VARCHAR(128)
	);
	`
	_, err := db.Exec(createTableSQL)
	if err != nil {
		return err
	}

	createIndexSQL := `CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler(date);`
	_, err = db.Exec(createIndexSQL)
	return err
}