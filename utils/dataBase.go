package utils

import (
	"context"
	"log"
	"os"

	db "main/sql"

	"github.com/jackc/pgx/v5"
)

func InitDatabase() (*db.Queries, *pgx.Conn) {
	// connect to db
	dsn := os.Getenv("dbType") + "://" + os.Getenv("dbUser") + ":" + os.Getenv("dbPwd") + "@" + os.Getenv("dbIP") + ":" + os.Getenv("dbPort") + "/" + os.Getenv("dbName") + "?sslmode=allow"

	conn, err := pgx.Connect(context.Background(), dsn)

	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close(context.Background())

	queries := db.New(conn)

	return queries, conn
}
