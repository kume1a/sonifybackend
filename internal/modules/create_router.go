package modules

import (
	"github.com/go-chi/cors"
	"github.com/gorilla/mux"
	"github.com/kume1a/sonifybackend/internal/modules/auth"
	"github.com/kume1a/sonifybackend/internal/modules/user"
	"github.com/kume1a/sonifybackend/internal/modules/youtubemusic"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func CreateRouter(apiCfg *shared.ApiConfg) *mux.Router {
	router := mux.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := router.PathPrefix("/v1").Subrouter().StrictSlash(false)

	v1Router.Handle("", auth.Router(apiCfg, v1Router))
	v1Router.Handle("", user.Router(apiCfg, v1Router))
	v1Router.Handle("", youtubemusic.Router(apiCfg, v1Router))

	router.Handle("", v1Router)
	router.HandleFunc("/", HandlerHealthcheck).Methods("GET")

	return router
}
