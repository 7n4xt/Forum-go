package main

import (
	"fmt"
	"forum-go/config"
	"forum-go/routes"
	"forum-go/templates"
	"log"
	"net/http"
)

func main() {
	config.LoadEnv()          // Load environment variables
	config.InitDB()           // Initialize database connection
	templates.LoadTemplates() // Load HTML templates and public files


	// Create a new ServeMux
	mux := http.NewServeMux()

	// Register routes
	routes.AuthRouter(mux)
	routes.DiscussionRouter(mux)
	routes.MessageRouter(mux)

	// Index page (home)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, "/discussions", http.StatusSeeOther)
	})

	// Set up the server
	port := "8080" // You can change this or make it configurable
	fmt.Printf("Server starting on port %s...\n", port)

	// Use the mux when starting the server
	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
