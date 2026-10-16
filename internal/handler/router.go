package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(baseURL string, db URLStorage) http.Handler {
	router := chi.NewRouter()
	router.Post("/", MakeRegisterURL(baseURL, db))
	router.Get("/{id}", MakeGetURL(db))
	return router
}
