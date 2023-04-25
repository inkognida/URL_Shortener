package data

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"math/rand"
	"shortener/pkg/utils"
	"time"
)

type Data struct {
	urls   map[string]string
	urlsID map[int64]struct{}
	logger *logrus.Logger
}

func (d *Data) Init(ctx context.Context, logger *logrus.Logger) error {
	d.urls = make(map[string]string)
	d.urlsID = make(map[int64]struct{})
	d.logger = logger
	return nil
}

func (d *Data) GetUrl(ctx context.Context, shortUrl string) (string, error) {
	if url, ok := d.urls[shortUrl]; ok {
		return url, nil
	}
	return "", errors.New("no such link")
}

func (d *Data) SaveUrl(ctx context.Context, originalUrl string) (string, error) {
	rand.Seed(time.Now().UnixNano())
	for {
		id := rand.Int63()
		if _, ok := d.urlsID[id]; !ok {
			d.urlsID[id] = struct{}{}
			shortUrl := utils.GenerateLink(id)
			d.urls[shortUrl] = originalUrl
			return shortUrl, nil
		}
	}
}
