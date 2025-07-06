package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/glebarez/sqlite"
)

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("sqlite", "./files.db")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS files (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		size INTEGER,
		uploaded_at TEXT
	);`
	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}

	fmt.Println("Database initialized ✅")
}

func saveFileToDB(name string, size int64) error {
	uploadedAt := time.Now().Format("2006-01-02 15:04:05")
	_, err := db.Exec("INSERT INTO files (name, size, uploaded_at) VALUES (?, ?, ?)", name, size, uploadedAt)
	return err
}
