package main

import "github.com/go-chi/chi/v5"

func (app *application) routes() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/v1/healthcheck", app.healthcheckHandler)

	return router
}
