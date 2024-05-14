package db

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

var DB *pgxpool.Pool

func InitDB() {
	// подключение к базе в формате postgres://myuser:mypassword@localhost:5432/mydatabase
	connString := os.Getenv("postgres://papawfen:@localhost:5432/papawfen") 
	var err error
	DB, err = pgxpool.Connect(context.Background(), connString)
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}

	createTable()
}

func createTable() {
	_, err := DB.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS commands (
			id UUID PRIMARY KEY,
			command TEXT NOT NULL,
			output TEXT,
			status TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatal("Unable to create table:", err)
	}
}
