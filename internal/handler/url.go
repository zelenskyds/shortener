package handler

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

func MakeRegisterURL(baseURL string, db URLStorage) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		data, err := io.ReadAll(req.Body)

		if err != nil {
			http.Error(res, "Не удалось прочитать body", http.StatusBadRequest)
			return
		}

		id := db.Save(string(data))

		body := fmt.Sprintf("%s/%s", strings.TrimRight(baseURL, "/"), id)

		res.Header().Set("Content-Type", "text/plain")
		res.WriteHeader(http.StatusCreated)
		res.Write([]byte(body))
	})
}

func MakeGetURL(db URLStorage) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		id := chi.URLParam(req, "id")

		url, ok := db.Get(id)

		if !ok {
			http.Error(res, "Url отсутствует", http.StatusNotFound)
			return
		}

		http.Redirect(res, req, url, http.StatusTemporaryRedirect)
	})
}
