package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"log"
	"net/http"
)

func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()

	// specify who is allowed to connect
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Use(middleware.Heartbeat("/ping"))
	mux.Use(middleware.DefaultLogger)

	mux.Post("/authenticate", app.exceptionMiddleware(app.Authenticate))

	return mux
}

type ExceptionMiddleware func(w http.ResponseWriter, r *http.Request) error

func (app *Config) exceptionMiddleware(f ExceptionMiddleware) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			if ierr := app.errorJSON(w, err, http.StatusBadRequest); ierr != nil {
				log.Println("cloud not send json:", err.Error())
			}

			return
		}

	}
}
