package main

import (
	"log"

	"github.com/joho/godotenv"
	"tomasweigenast.com/url-shortener/server"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error loading .env file: %s", err)
	}
}

func main() {
	server.Run()
}
