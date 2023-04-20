package data

import (
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

func (d *Data) Init() {
	d.urls = make(map[string]string)
	d.urlsID = make(map[int64]struct{})
}

func (d *Data) GetUrl(shortUrl string) (string, error) {
	if url, ok := d.urls[shortUrl]; ok {
		return url, nil
	}
	return "", errors.New("url not found")
}

func (d *Data) SaveUrl(originalUrl string) (string, error) {
	rand.Seed(time.Now().UnixNano())
	// TODO fix forever running
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
