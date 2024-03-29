package main

import (
	"context"
	"fmt"
	"main/router"
	db "main/sql"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	// connect to db
	dsn := os.Getenv("dbType") + "://" + os.Getenv("dbUser") + ":" + os.Getenv("dbPwd") + "@" + "dbIP" + ":" + os.Getenv("dbPort") + "/" + os.Getenv("dbName") + "?sslmode=allow"

	conn, err := pgx.Connect(context.Background(), dsn)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	queries := db.New(conn)

	router.InitRouter()
}
