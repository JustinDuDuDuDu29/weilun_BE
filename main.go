package main

import (
	"main/router"
	"main/utils"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	q, conn := utils.InitDatabase()

	ctrl := utils.Init(q, conn)
	defer conn.Close()
	router.RouterInit(ctrl)
}
