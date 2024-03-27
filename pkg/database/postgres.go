package database

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/ccthomas/gridiron/pkg/logger"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func ConnectPostgres() *sql.DB {
	logger.Get().Debug("Connecting to postgres")

	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Do not log connection string, as it contains sensitive information.
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:5432/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbName)

	logger.Get().Debug("Open a connection to the postgres database")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logger.Get().Fatal("Error opening database connection:", zap.Error(err))
	}

	logger.Get().Debug("Connection to postgres database was successful.")
	return db
}
