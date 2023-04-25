package server

import (
	"context"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net"
	"os"
	"os/signal"
	shortener "shortener/pkg/api/grpc"
	repo "shortener/pkg/repository"
	"syscall"
)

type server struct {
	shortener.UnsafeUrlShortenerServer
	storage repo.StorageType
	logger *logrus.Logger
}

// ShortenUrl метод для сокращения оригинальной ссылки
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

// ExtractUrl метод для получения оригинальной ссылки по сокращенной
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

// InitServer функция инициализирует сервер с параметрамом для храналища
func InitServer(logger *logrus.Logger, storageMode string) shortener.UrlShortenerServer {
	srv := &server{logger: logger}
	storage, err := repo.NewRepo(storageMode, logger)
	if err != nil {
		srv.logger.Fatalln("failed to create repo", err)
	}
	srv.storage = storage
	return srv
}

// Execute создает и запускает новый grpc сервер
func Execute(logger *logrus.Logger, storageMode string) {
	srv := grpc.NewServer()
	shortener.RegisterUrlShortenerServer(srv, InitServer(logger, storageMode))

	listen, err := net.Listen("tcp", ":8000")
	if err != nil {
		logger.Log(logrus.FatalLevel, err)
	}
	errC := make(chan error, 1)
	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-ctx.Done()
		logger.Log(logrus.DebugLevel, "got signal, starting graceful shutdown")
		defer func() {
			stop()
			close(errC)
		}()
		srv.GracefulStop()
	}()

	logger.Infoln("starting grpc server")
	go func() {
		if err = srv.Serve(listen); err != nil {
			errC <- err
		}
	}()

	if err = <-errC; err != nil {
		logger.Log(logrus.DebugLevel, "could not serve", err)
	}

	if err = <-errC; err != nil {
		logger.Fatalln("Error while execution", err)
	}
}

