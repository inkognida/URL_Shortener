package repo

import (
	"github.com/sirupsen/logrus"
	"shortener/pkg/repository/data"
)

type StorageType interface {
	Init()
	GetUrl(shortUrl string) (string, error)
	SaveUrl(originalUrl string) (string, error)
}

func NewRepoStorage(mode string, logger *logrus.Logger) StorageType {
	if mode == "postgres" {
		// TODO implement postgres storage
	}
	storage := &data.Data{}
	storage.Init()
	return storage
}
