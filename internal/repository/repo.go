package repo

import (
	"context"
	"github.com/sirupsen/logrus"
	"shortener/internal/repository/data"
	"shortener/internal/repository/postgres"
)

// StorageType интерфейс хранилища
type StorageType interface {
	// Init инициализация хранилища
	Init(ctx context.Context, logger *logrus.Logger) error

	// GetUrl возвращает оригинальную ссылку по сокращенной
	GetUrl(ctx context.Context, shortUrl string) (string, error)

	// SaveUrl сохрняет сокращенный вариант ссылки и возвращает ее
	SaveUrl(ctx context.Context, originalUrl string) (string, error)
}

// NewRepo функция для создания нового репризитория под хранилище с его параметром
func NewRepo(mode string, logger *logrus.Logger) (StorageType, error) {
	if mode == "postgres" {
		storage := &postgres.Data{}
		err := storage.Init(context.Background(), logger)
		if err != nil {
			return nil, err
		}

		logger.Infoln("repo created with mode:", mode)
		return storage, nil
	}

	storage := &data.Data{}
	if err := storage.Init(context.Background(), logger); err != nil {
		return nil, err
	}

	logger.Infoln("repo created with mode:", mode)
	return storage, nil
}
