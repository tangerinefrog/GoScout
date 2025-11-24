package main

import (
	"job-scraper/internal/data"
	"job-scraper/internal/handlers"
	"log"

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

	port := "8080"
	log.Printf("Server is starting on port :%s...\n", port)
	err = r.Run(":" + port)
	log.Printf("Server error: %v\n", err)
}
