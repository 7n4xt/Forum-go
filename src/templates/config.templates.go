package templates

import (
	"html/template"
	"log"
)

var Temp *template.Template

func LoadTemplates() {
	listeTemplates, errTemplates := template.ParseGlob("./templates/*.html")
	if errTemplates != nil {
		log.Fatalf("Erreur chargement templates - Erreur : %s", errTemplates.Error())
	}
	Temp = listeTemplates
}
