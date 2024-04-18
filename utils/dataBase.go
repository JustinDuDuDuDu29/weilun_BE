package utils

import (
	"database/sql"
	"log"
	"os"

	db "main/sql"

	_ "github.com/lib/pq"
)

func InitDatabase() (*db.Queries, *sql.DB) {
	// connect to db
	dsn := os.Getenv("dbType") + "://" + os.Getenv("dbUser") + ":" + os.Getenv("dbPwd") + "@" + os.Getenv("dbIP") + ":" + os.Getenv("dbPort") + "/" + os.Getenv("dbName") + "?sslmode=disable"

	conn, err := sql.Open("postgres", dsn)

	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	// defer conn.Close(context.Background())

	queries := db.New(conn)

	return queries, conn
}
