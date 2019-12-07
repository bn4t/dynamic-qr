package main

import (
	"git.bn4t.me/bn4t/dynamic-qr/app/utils"
	"git.bn4t.me/bn4t/dynamic-qr/app/web"
	"git.bn4t.me/bn4t/dynamic-qr/db"
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

	// open the database
	db.Connect()

	// start webserver
	web.Start()

}
