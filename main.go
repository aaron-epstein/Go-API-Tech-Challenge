package main

import (
	"log"

	"github.com/aaron-epstein/Go-API-Tech-Challenge/internal"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	_, err = internal.InitDB()
	if err != nil {
		// panic(err)
		log.Fatal("Error connecting to DB")
	}

	internal.RunServer()

}
