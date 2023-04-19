package main

import (
	"errors"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"net"
	shortener "shortener/internal/shortenerpb"
	"strings"
	"time"
)


type StorageType interface {
	Init()
	GetUrl(shortUrl string)
	SaveUrl(originalUrl string)
}

type server struct {
	shortener.UnsafeUrlShortenerServer
	storage StorageType
}

type Data struct {
	urls map[string]string
}

func (d *Data) Init() {
	d.urls = make(map[string]string, 0)
}

func (d *Data) GetUrl(shortUrl string) (string, error){
	if url, ok := d.urls[shortUrl]; ok {
		return url, nil
	}
	return "", errors.New("url not found")
}

func (d *Data) Save(originalUrl string) (string, error) {
	rand.Seed(time.Now().UnixNano())
	id := rand.Int63()

	for {
		link := GenerateLink(id)
		if _, ok := d.urls[link]; !ok {
			d.urls[link] = originalUrl
			return link, nil
		}
		id = rand.Int63()
	}
}

const (
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	linkLength = 10
	charsetLength = int64(len(charset))
)

func GenerateLink(shortId int64) string {
	if shortId < 0 {
		shortId = -shortId
	}
	var b [11]byte
	for i := linkLength; shortId > 0 && i > 0; i-- {
		shortId, b[i] = shortId/charsetLength, charset[shortId%charsetLength]
	}
	return string(b[:])
}

func main() {
	logger := logrus.New()

	listen, err := net.Listen("tcp", ":8000")
	if err != nil {
		logger.Log(logrus.FatalLevel, err)
	}

	srv := grpc.NewServer()
	// TODO CHECK
	defer srv.GracefulStop()

	// TODO add PostgreSQL realization
	shortener.RegisterUrlShortenerServer(srv, &server{
		storage: ,
	})

	candlesStorage, err := candles.LoadFromFile("candles",
		[]domain.CandlePeriod{domain.CandlePeriod1m, domain.CandlePeriod2m, domain.CandlePeriod10m})
	if err != nil {
		log.Fatalf("can't load candles: %v", err)
	}

	candlespb.RegisterCandlesServiceServer(s, &server{candles: candlesStorage})
	if err := s.Serve(listen); err != nil {
		log.Fatalf("can't register service server: %v", err)
	}
}
