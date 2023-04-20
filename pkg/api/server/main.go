package main

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
}

func (s *server) ShortenUrl(ctx context.Context, request *shortener.ShortenURLRequest) (*shortener.ShortenURLResponse, error) {
	if request.OriginUrl == "" {
		return nil, status.Errorf(codes.NotFound, "empty link")
	}
	link, err := s.storage.SaveUrl(request.OriginUrl)
	if err != nil {
		return nil, err
	}

	return &shortener.ShortenURLResponse{ShortenedUrl: link}, nil
}

func (s *server) ExtractUrl(ctx context.Context, request *shortener.ExtractURLRequest) (*shortener.ExtractURLResponse, error) {
	link, err := s.storage.GetUrl(request.ShortenedUrl)
	if err != nil {
		return nil, err
	}

	return &shortener.ExtractURLResponse{OriginUrl: link}, nil
}

func main() {
	logger := logrus.New()

	// TODO add config
	listen, err := net.Listen("tcp", ":8000")
	if err != nil {
		logger.Log(logrus.FatalLevel, err)
	}
	srv := grpc.NewServer()
	// TODO add PostgreSQL realization
	shortener.RegisterUrlShortenerServer(srv, &server{
		storage: repo.NewRepo("data", logger),
	})

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

	logrus.Infoln("starting grpc server")
	go func() {
		if err = srv.Serve(listen); err != nil {
			errC <- err
		}
	}()

	if err = <-errC; err != nil {
		logger.Log(logrus.DebugLevel, "could not serve")
	}

	logrus.Infoln("clean shutdown")
}
