package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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

	workEnqueuer, workPool := bgwork.ConfigureBackgroundWork(resouceConfig)

	apiCfg := config.ApiConfig{
		ResourceConfig: resouceConfig,
		WorkEnqueuer:   workEnqueuer,
	}

	router := modules.CreateRouter(&apiCfg)

	server := &http.Server{
		Handler: router,
		Addr:    ":" + envVars.Port,
	}

	log.Println("Server is running on port:", envVars.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	<-signalChan
	workPool.Stop()
}
