package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/modules"
	"github.com/kume1a/sonifybackend/internal/modules/bgwork"
	"github.com/kume1a/sonifybackend/internal/modules/ws"
	"github.com/kume1a/sonifybackend/internal/shared"

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
	router.PathPrefix("/").Handler(http.StripPrefix("/public", http.FileServer(http.Dir("public/"))))
	router.HandleFunc("/ws", ws.HandleWsUpgrade)

	server := &http.Server{
		Handler: router,
		Addr:    ":" + envVars.Port,
	}

	go shared.Ticker(1*time.Second, 100*time.Second, func() {
		ws, exists := ws.GetManager().GetConnection("uniqueKey")
		if !exists {
			return
		}

		// Use conn to write messages to the client
		if err := ws.WriteMessage(websocket.TextMessage, []byte("Your message here")); err != nil {
			log.Println("Error writing to websocket:", err)
			return
		}
	})

	log.Println("Server is running on port:", envVars.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	<-signalChan
	workPool.Stop()
}
