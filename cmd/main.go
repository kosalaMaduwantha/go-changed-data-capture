package main

import (
	"cdc-file-processor/domain"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	domain.Cdc_run()
}