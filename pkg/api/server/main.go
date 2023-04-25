package server

import (
	"context"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	shortener "shortener/pkg/api/grpc"
	repo "shortener/pkg/repository"
)

type server struct {
	shortener.UnsafeUrlShortenerServer
	storage repo.StorageType
	logger *logrus.Logger
}

func (s *server) ShortenUrl(ctx context.Context, request *shortener.ShortenURLRequest) (*shortener.ShortenURLResponse, error) {
	if request.OriginUrl == "" {
		s.logger.Infoln("empty link via request:", request.OriginUrl)
		return &shortener.ShortenURLResponse{ShortenedUrl: request.OriginUrl}, status.Errorf(codes.NotFound, "empty link")
	}

	link, err := s.storage.SaveUrl(ctx, request.OriginUrl)
	if err != nil {
		s.logger.Infoln("failed to save url:", err)
		return &shortener.ShortenURLResponse{ShortenedUrl: link}, err
	}

	return &shortener.ShortenURLResponse{ShortenedUrl: link}, nil
}

func (s *server) ExtractUrl(ctx context.Context, request *shortener.ExtractURLRequest) (*shortener.ExtractURLResponse, error) {
	if request.ShortenedUrl == "" || len(request.ShortenedUrl) != 10 {
		s.logger.Infoln("wrong link parameter via request:", request.ShortenedUrl)
		return &shortener.ExtractURLResponse{OriginUrl: request.ShortenedUrl}, status.Errorf(codes.Unavailable, "wrong link")
	}

	link, err := s.storage.GetUrl(ctx, request.ShortenedUrl)
	if err != nil {
		s.logger.Infoln("failed to get url:", err)
		return &shortener.ExtractURLResponse{OriginUrl: link}, err
	}

	return &shortener.ExtractURLResponse{OriginUrl: link}, nil
}

func InitServer(storageMode string, logger *logrus.Logger) shortener.UrlShortenerServer {
	srv := &server{logger: logger}
	storage, err := repo.NewRepo(storageMode, logger)
	if err != nil {
		srv.logger.Fatalln("failed to create repo", err)
	}
	srv.storage = storage
	return srv
}

