package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/modules"
	"github.com/kume1a/sonifybackend/internal/shared"

	_ "github.com/lib/pq"
)

// check lib for http: https://github.com/levigross/grequests?utm_campaign=awesomego&utm_medium=referral&utm_source=awesomego

// TODOs
// - update playlist changes in db when reimporting from spotify
// - import all audios and all playlist, use paging of spotify limit is currently 50

func main() {
	shared.LoadEnv()

	envVars, err := shared.ParseEnv()
	if err != nil {
		log.Fatal("Coultn't parse env vars, returning")
		return
	}

	govalidator.SetFieldsRequiredByDefault(true)

	conn, err := sql.Open("postgres", envVars.DbUrl)
	if err != nil {
		log.Fatal("Couldn't connect to database", envVars.DbUrl)
		return
	}

	apiCfg := shared.ApiConfig{
		DB:    database.New(conn),
		SqlDB: conn,
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
