package db

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strconv"
	"time"
)

func Connect() (*sql.DB, error) {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	var DB *sql.DB
	var err error

	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return nil, err
	}
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		port,
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"))

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	//// Выполняем миграции с помощью goose
	//migrationDir := "./db/migrations"
	//if err := goose.Up(db, migrationDir); err != nil {
	//	log.Fatalf("[ERROR] Could not run migrations: %v", err)
	//	return err
	//}

	start := time.Now()

	if err = DB.Ping(); err != nil {
		log.Printf("[ERROR] Database ping failed: %v", err)
		return nil, err
	}

	elapsed := time.Since(start)
	log.Printf("[INFO] Database connection is successful. Ping time: %s", elapsed)

	log.Println("[INFO] Database connection established successfully")
	return DB, nil
}
