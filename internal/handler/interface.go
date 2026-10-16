package handler

type URLStorage interface {
	Save(url string) string
	Get(id string) (string, bool)
}
