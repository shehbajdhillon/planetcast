package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"

	"go.uber.org/zap"
)

type DatabaseConnectProps struct {
	Logger *zap.Logger
}

func Connect(args DatabaseConnectProps) *Queries {

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	dbname := os.Getenv("DB_NAME")

	logger := args.Logger

	postgresqlDbInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	db, err := sql.Open("postgres", postgresqlDbInfo)

	if err != nil {
		logger.Fatal("Could not connect to database", zap.Error(err))
	}

	err = db.Ping()

	if err != nil {
		logger.Fatal("Could not connect to database", zap.Error(err))
	}

	logger.Info("Database client started")

	return New(db)
}
