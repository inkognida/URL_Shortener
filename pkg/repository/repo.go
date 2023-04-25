package repo

import (
	"context"
	"github.com/sirupsen/logrus"
	"shortener/pkg/repository/data"
	"shortener/pkg/repository/postgres"
)

type StorageType interface {
	Init(ctx context.Context, logger *logrus.Logger) error
	GetUrl(ctx context.Context, shortUrl string) (string, error)
	SaveUrl(ctx context.Context, originalUrl string) (string, error)
}

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
