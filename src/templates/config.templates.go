package templates

import (
	"html/template"
	"log"
	"net/http"
)

var Temp *template.Template

func LoadTemplates() {
	templates := []string{
		"./templates/*.html",
		"./templates/*.template.html",
		"./templates/*.templates.html",
	}

	// Configure static file server for CSS, JS, and images
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/public/", http.StripPrefix("/public/", fs))

	var err error
	Temp = template.New("")

	for _, pattern := range templates {
		Temp, err = Temp.ParseGlob(pattern)
		if err != nil {
			log.Printf("Warning: Error loading templates from %s: %v", pattern, err)
		}
	}

	if Temp == nil {
		log.Fatalf("Error: No templates were loaded")
	}
}
