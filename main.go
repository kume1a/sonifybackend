package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	envVars, err := parseEnv()
	if err != nil {
		log.Fatal("Coultn't parse env vars, returning")
		return
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
	v1Router.Get("/healthcheck", healthcheck)
	v1Router.Get("/err", handleErr)

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
