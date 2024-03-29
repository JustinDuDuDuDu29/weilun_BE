package main

import (
	"main/router"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	router.InitRouter()
}
