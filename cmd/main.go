package main

import (
	"job-scraper/internal/data"
	"job-scraper/internal/handlers"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("could not load .env file: %v", err)
	}

	db, err := data.Init()
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	h := handlers.NewHandler(db)
	h.SetupRoutes(r)

	addr := os.Getenv("SRV_ADDR")
	if addr == "" {
		log.Fatalf("Server address is not defined in the .env file")
	}

	log.Printf("Server is listening on ':%s'...\n", addr)
	err = r.Run(addr)
	log.Printf("Server error: %v\n", err)
}
