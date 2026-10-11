package handler

import (
	"crypto/rand"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type DB map[string]string

func MakeRegistrURL(baseURL string, baseGetURL string, db DB) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		data, err := io.ReadAll(req.Body)

		if err != nil {
			http.Error(res, "Не удалось прочитать body", http.StatusBadRequest)
			return
		}

		id := rand.Text()

		db[id] = string(data)

		body := fmt.Sprintf("http://%s/%s", baseGetURL, id)

		res.Header().Set("Content-Type", "text/plain")
		res.WriteHeader(http.StatusCreated)
		res.Write([]byte(body))
	})
}

func MakeGetURL(db DB) http.HandlerFunc {
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

func newRouter(baseURL string, baseGetURL string, db DB) http.Handler {
	router := chi.NewRouter()
	router.Post("/", MakeRegistrURL(baseURL, baseGetURL, db))
	router.Get("/{id}", MakeGetURL(db))
	return router
}

func RunApp(baseURL string, baseGetURL string) {
	db := make(DB)
	err := http.ListenAndServe(baseURL, newRouter(baseURL, baseGetURL, db))
	if err != nil {
		panic(err)
	}
}
