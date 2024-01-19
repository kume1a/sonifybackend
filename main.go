package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/cors"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/modules/user"
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

	conn, err := sql.Open("postgres", envVars.DbUrl)
	if err != nil {
		log.Fatal("Couldn't connect to database", envVars.DbUrl)
	}

	apiCfg := shared.ApiConfg{
		DB: database.New(conn),
	}

	router := mux.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := mux.NewRouter()
	// v1Router.Handle("/healthcheck", modules.HandlerHealthcheck).Methods("GET")

	v1Router.Handle("/users", user.HandlerCreateUser(&apiCfg)).Methods("POST")

	router.Handle("/v1", v1Router)

	server := &http.Server{
		Handler: router,
		Addr:    ":" + envVars.Port,
	}

	log.Printf("Starting server on port %s", envVars.Port)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
