package main

import (
	"main/router"
	"main/utils"
	"main/wiregen"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	q, conn := utils.InitDatabase()

	ctrl := wiregen.Init(q, conn)
	mid := wiregen.MInit(q, conn)
	defer conn.Close()
	router.RouterInit(ctrl, mid)
}
