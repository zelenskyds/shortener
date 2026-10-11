package handler

import (
	"crypto/rand"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type DB map[string]string

func MakeRegistrUrl(base_url string, db DB) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		data, err := io.ReadAll(req.Body)

		if err != nil {
			http.Error(res, "Не удалось прочитать body", http.StatusBadRequest)
			return
		}

		id := rand.Text()

		db[id] = string(data)

		body := fmt.Sprintf("http://%s/%s", base_url, id)

		res.Header().Set("Content-Type", "text/plain")
		res.WriteHeader(http.StatusCreated)
		res.Write([]byte(body))
	})
}

func MakeGetUrl(db DB) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		id := chi.URLParam(req, "id")

		url, ok := db[id]

		if !ok {
			http.Error(res, "Url отсутствует", http.StatusNotFound)
			return
		}

		http.Redirect(res, req, url, http.StatusTemporaryRedirect)
	})
}

func newRouter(base_url string, db DB) http.Handler {
	router := chi.NewRouter()
	router.Post("/", MakeRegistrUrl(base_url, db))
	router.Get("/{id}", MakeGetUrl(db))
	return router
}

func RunApp(base_url string) {
	db := make(DB)
	err := http.ListenAndServe(base_url, newRouter(base_url, db))
	if err != nil {
		panic(err)
	}
}
