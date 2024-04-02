package main

import (
	"context"
	"fmt"
	"main/router"
	"main/utils"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	fmt.Print(os.Getenv("dbUser"))
	q, conn := utils.InitDatabase()

	ctrl := utils.Init(q)
	defer conn.Close(context.Background())
	router.RouterInit(ctrl)
}
