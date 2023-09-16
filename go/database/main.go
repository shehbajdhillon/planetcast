package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func Connect() *Queries {

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	dbname := os.Getenv("DB_NAME")

	postgresqlDbInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	db, err := sql.Open("postgres", postgresqlDbInfo)

	if err != nil {
		log.Fatalln("Could not connect to database:", err.Error())
	}

	err = db.Ping()

	if err != nil {
		log.Fatalln("Could not connect to database:", err.Error())
	}

	log.Println("Database client started")

	return New(db)
}
