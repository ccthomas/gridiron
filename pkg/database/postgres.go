package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func (dm *DatabaseManager) ConnectPostgres() *sql.DB {
	dm.Logger.Debug("Connecting to postgres")

	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Do not log connection string, as it contains sensitive information.
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:5432/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbName)

	dm.Logger.Debug("Open a connection to the postgres database")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error opening database connection:", err)
	}

	dm.Logger.Debug("Connection to postgres database was successful.")
	return db
}
