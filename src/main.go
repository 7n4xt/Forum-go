package main

import (
	"fmt"
	"forum-go/config"
	"log"
	"net/http"
	"time"
)

func main() {
	config.LoadEnv() // Load environment variables
 config.InitDB() // Initialize database connection
// Configure static file server for CSS, JS, and images
fs := http.FileServer(http.Dir("public"))
http.Handle("/public/", http.StripPrefix("/public/", fs))

// Handle the main route
http.HandleFunc("/", indexHandler)

// Server configuration
server := &http.Server{
Addr:         ":8080",
ReadTimeout:  10 * time.Second,
WriteTimeout: 10 * time.Second,
}

fmt.Println("Starting server at http://localhost:8080")
log.Fatal(server.ListenAndServe())
}

// Handle the index route
func indexHandler(w http.ResponseWriter, r *http.Request) {
http.ServeFile(w, r, "templates/index.html")
}
