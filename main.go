package main

import (
	"fmt"
	"main/router"
	"main/utils"
	"main/wiregen"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	fmt.Println("Starting version: ", os.Getenv("version"))

	godotenv.Load()
	q, conn := utils.InitDatabase()

	ctrl := wiregen.Init(q, conn)
	mid := wiregen.MInit(q, conn)
	defer conn.Close()
	router.RouterInit(ctrl, mid)
}
