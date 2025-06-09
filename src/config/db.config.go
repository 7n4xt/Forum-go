package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DbContext *sql.DB

func InitDB() {

	user := GetEnvWithDefault("DB_USER", "")
	pwd := GetEnvWithDefault("DB_PWD", "")
	host := GetEnvWithDefault("DB_HOST", "localhost")
	port := GetEnvWithDefault("DB_PORT", "3306")
	name := GetEnvWithDefault("DB_NAME", "")

	if user == "" || name == "" {
		log.Fatalf("Erreur de connexion à la base de données : paramètres manquants")
	}

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, pwd, host, port, name)

	var err error
	DbContext, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatalf("Erreur d'ouverture de la connexion : %v", err)
	}

	err = DbContext.Ping()
	if err != nil {

		DbContext.Close()
		log.Fatalf("Échec du ping vers MySQL : %v", err)
	}
}
