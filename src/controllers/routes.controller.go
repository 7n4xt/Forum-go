package controllers

import (
	"net/http"
)

// SetupRoutes configures all application routes
func SetupRoutes() {
	// Static files
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/public/", http.StripPrefix("/public/", fs))

	// Main routes
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/register", RegisterHandler)
	http.HandleFunc("/login", LoginHandler)
	http.HandleFunc("/logout", LogoutHandler)
}

// IndexHandler handles the main page
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, "templates/index.html")
}
