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
	// ticker := time.NewTicker(5000 * time.Millisecond)
	// go func() {
	// 	for range ticker.C {
	// 		res, err := http.Get("https://www.imdu29.com")
	// 		fmt.Println(res.StatusCode)
	// 		if res.StatusCode == 404 {
	// 			fmt.Println("404!")
	// 			ticker.Stop()
	// 			return
	// 		}
	// 		if res.StatusCode != 200 || err != nil {
	// 			fmt.Println("not 200!")
	// 			os.Exit(1)
	// 		}
	// 	}
	// }()
	fmt.Println("Starting version: ", os.Getenv("version"))

	godotenv.Load()
	q, conn := utils.InitDatabase()

	ctrl := wiregen.Init(q, conn)
	mid := wiregen.MInit(q, conn)
	defer conn.Close()
	router.RouterInit(ctrl, mid)
}
