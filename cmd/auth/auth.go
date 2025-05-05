package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"iztech-agms/router"
	"log"
	"os"
	"path/filepath"
)

func main() {
	pwd, _ := os.Getwd()
	environmentPath := filepath.Join(pwd, ".env")
	err := godotenv.Load(environmentPath)
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	blackListHosts := map[string]struct{}{}

	// initialize new router and add handlers
	r := router.AuthRouter(blackListHosts)

	port := os.Getenv("PORT")
	port = fmt.Sprintf(":%s", port)

	err = r.Run(port)
	if err != nil {
		log.Fatal(err)
	}
}
