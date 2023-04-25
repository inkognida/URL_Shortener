package handler

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	repo "shortener/internal/repository"
)

var storage repo.StorageType

// NewHandler функция для создания нового обработчика с параметром для хранилища
func NewHandler(storageMode string, logger *logrus.Logger) (chi.Router, error) {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Post("/", SaveUrl)
	router.Get("/{shortUrl}", GetUrl)

	var err error
	storage, err = repo.NewRepo(storageMode, logger)
	if err != nil {
		return nil, err
	}

	return router, nil
}

// SaveUrl сокращение оригинальной ссылки
func SaveUrl(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	originalURL := string(bodyBytes)
	shortUrl, err := storage.SaveUrl(context.Background(), originalURL)
	if err != nil {
		w.WriteHeader(500)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(200)
	_, _ = w.Write([]byte(shortUrl))
}

// GetUrl поулчение оригинальной ссылки по сокращенной
func GetUrl(w http.ResponseWriter, r *http.Request) {
	shortUrl := chi.URLParam(r, "shortUrl")
	originalURL, err := storage.GetUrl(context.Background(), shortUrl)
	if err != nil {
		w.WriteHeader(404)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(200)
	_, _ = w.Write([]byte(originalURL))
}
