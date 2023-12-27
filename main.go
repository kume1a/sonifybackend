package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/kume1a/sonifybackend/internal/database"

	_ "github.com/lib/pq"
)

type apiConfg struct {
	DB *database.Queries
}

func main() {
	godotenv.Load(".env")

	envVars, err := parseEnv()
	if err != nil {
		log.Fatal("Coultn't parse env vars, returning")
		return
	}

	conn, err := sql.Open("postgres", envVars.DbUrl)
	if err != nil {
		log.Fatal("Couldn't connect to database", envVars.DbUrl)
	}

	apiCfg := apiConfg{
		DB: database.New(conn),
	}

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/healthcheck", handlerHealthcheck)
	v1Router.Get("/err", handlerErr)

	v1Router.Post("/users", apiCfg.handlerCreateUser)

	router.Mount("/v1", v1Router)

	server := &http.Server{
		Handler: router,
		Addr:    ":" + envVars.Port,
	}

	log.Printf("Starting server on port %s", envVars.Port)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
