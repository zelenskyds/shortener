package handler

import (
	"crypto/rand"
	"fmt"
	"io"
	"net/http"
)

type DB map[string]string

func MakeRegistrUrl(base_url string, db DB) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			http.Error(res, "Метод неразрешён", http.StatusMethodNotAllowed)
			return
		}

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
		if req.Method != http.MethodGet {
			http.Error(res, "Метод неразрешён", http.StatusMethodNotAllowed)

			return
		}

		id := req.PathValue("id")

		url, ok := db[id]

		if !ok {
			http.Error(res, "Url отсутствует", http.StatusNotFound)
			return
		}

		http.Redirect(res, req, url, http.StatusTemporaryRedirect)
	})
}

func RunApp(base_url string) {
	db := make(DB)

	mux := http.NewServeMux()

	mux.Handle(`/{$}`, MakeRegistrUrl(base_url, db))
	mux.Handle(`/{id}`, MakeGetUrl(db))

	err := http.ListenAndServe(base_url, mux)
	if err != nil {
		panic(err)
	}
}
