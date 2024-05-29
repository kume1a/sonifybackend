package modules

import (
	"errors"
	"fmt"

	"github.com/go-chi/cors"
	"github.com/gorilla/mux"
	"github.com/kume1a/sonifybackend/internal/config"
	"github.com/kume1a/sonifybackend/internal/modules/audio"
	"github.com/kume1a/sonifybackend/internal/modules/audiolike"
	"github.com/kume1a/sonifybackend/internal/modules/auth"
	"github.com/kume1a/sonifybackend/internal/modules/playlist"
	"github.com/kume1a/sonifybackend/internal/modules/spotify"
	"github.com/kume1a/sonifybackend/internal/modules/user"
	"github.com/kume1a/sonifybackend/internal/modules/userplaylist"
	"github.com/kume1a/sonifybackend/internal/modules/usersync"
	"github.com/kume1a/sonifybackend/internal/modules/youtube"
)

func CreateRouter(apiCfg *config.ApiConfig) *mux.Router {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)

			var err error
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("unknown panic")
			}
			if err != nil {
				// sendMeMail(err)
				fmt.Println("sendMeMail")
			}
		}
	}()

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
	v1Router.Handle("", youtube.Router(apiCfg, v1Router))
	v1Router.Handle("", audio.Router(apiCfg, v1Router))
	v1Router.Handle("", playlist.Router(apiCfg, v1Router))
	v1Router.Handle("", spotify.Router(apiCfg, v1Router))
	v1Router.Handle("", usersync.Router(apiCfg, v1Router))
	v1Router.Handle("", audiolike.Router(apiCfg, v1Router))
	v1Router.Handle("", userplaylist.Router(apiCfg, v1Router))

	router.Handle("", v1Router)

	router.HandleFunc("/", handleHealthcheck).Methods("GET")
	router.HandleFunc("/serverTime", handleGetServerTime).Methods("GET")

	return router
}
