package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/joho/godotenv"
	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/modules"
	"github.com/kume1a/sonifybackend/internal/shared"

	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load(".env")

	envVars, err := shared.ParseEnv()
	if err != nil {
		log.Fatal("Coultn't parse env vars, returning")
		return
	}

	govalidator.SetFieldsRequiredByDefault(true)

	conn, err := sql.Open("postgres", envVars.DbUrl)
	if err != nil {
		log.Fatal("Couldn't connect to database", envVars.DbUrl)
	}

	apiCfg := shared.ApiConfg{
		DB: database.New(conn),
	}

	router := modules.CreateRouter(&apiCfg)
	router.PathPrefix("/").Handler(http.StripPrefix("/public", http.FileServer(http.Dir("public/"))))

	server := &http.Server{
		Handler: router,
		Addr:    ":" + envVars.Port,
	}

	log.Printf("Starting server on port %s", envVars.Port)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
