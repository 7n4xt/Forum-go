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
	templates.LoadTemplates() // Load HTML templates
	// Configure static file server for CSS, JS, and images
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/public/", http.StripPrefix("/public/", fs))

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Register routes
	routes.AuthRouter(mux)

	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		c, cerr := r.Cookie("token")
		fmt.Println("Cookie:", c, "Error:", cerr)
		templates.Temp.ExecuteTemplate(w, "index", nil)
	})

	// You may have other routes to register
	// routes.OtherRouter(mux)

	// Set up the server
	port := "8080" // You can change this or make it configurable
	fmt.Printf("Server starting on port %s...\n", port)

	// Use the mux when starting the server
	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

	

}
