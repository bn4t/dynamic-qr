package main

import (
	"git.bn4t.me/bn4t/dynamic-qr/utils"
	"git.bn4t.me/bn4t/dynamic-qr/web"
	"github.com/joho/godotenv"
	"log"
)

func main() {

	if utils.FileExists(".env") {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		} else {
			log.Println("Successfully loaded environment variables")
		}
	}

	// start webserver
	web.Start()


}
