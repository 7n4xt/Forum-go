package main

import (
	"fmt"
	"forum-go/config"
	"forum-go/controllers"
	"log"
	"net/http"
	"time"
)

func main() {
	config.LoadEnv() // Load environment variables
	config.InitDB()  // Initialize database connection

	// Setup all routes
	controllers.SetupRoutes()

	// Server configuration
	server := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Println("Starting server at http://localhost:8080")
	log.Fatal(server.ListenAndServe())
}
