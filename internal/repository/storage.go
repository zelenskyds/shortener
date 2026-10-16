package repository

import (
	"crypto/rand"
)

var urlToIDDB = make(map[string]string)
var idToURLDB = make(map[string]string)

type Storage struct{}

func (r *Storage) Save(url string) string {
	oldID, ok := urlToIDDB[url]

	if ok {
		return oldID
	}

	id := rand.Text()

	idToURLDB[id] = url
	urlToIDDB[url] = id

	return id
}

func (r *Storage) Get(id string) (string, bool) {
	url, ok := idToURLDB[id]

	return url, ok
}
