package templates

import (
	"html/template"
	"log"
)

var Temp *template.Template

func LoadTemplates() {
	// Define template functions
	funcMap := template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"sub": func(a, b int) int {
			return a - b
		},
		"len": func(slice interface{}) int {
			switch s := slice.(type) {
			case []interface{}:
				return len(s)
			default:
				return 0
			}
		},
	}

	templates := []string{
		"./templates/*.html",
		"./templates/*.template.html",
		"./templates/*.templates.html",
	}

	var err error
	Temp = template.New("").Funcs(funcMap)

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
