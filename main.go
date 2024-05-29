package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/modules"
	"github.com/kume1a/sonifybackend/internal/modules/bgwork"

	"github.com/kume1a/sonifybackend/internal/config"

	_ "github.com/lib/pq"
)

func main() {
	config.LoadEnv()

	envVars, err := config.ParseEnv()
	if err != nil {
		log.Fatal("Coultn't parse env vars, returning")
		return
	}

	config.ConfigureGoValidator()

	conn, err := sql.Open("postgres", envVars.DbUrl)
	if err != nil {
		log.Fatal("Couldn't connect to database", envVars.DbUrl)
		return
	}

	db := database.New(conn)

	resouceConfig := &config.ResourceConfig{
		DB:    db,
		SqlDB: conn,
	}

	workEnqueuer := bgwork.ConfigureBackgroundWork(resouceConfig)

	apiCfg := config.ApiConfig{
		ResourceConfig: resouceConfig,
		WorkEnqueuer:   workEnqueuer,
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
