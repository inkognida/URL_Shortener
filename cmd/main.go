package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	shortener "shortener/pkg/api/grpc"
	"shortener/pkg/api/server"
	"syscall"
)

func execute(logger *logrus.Logger, srv *grpc.Server) {
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

	logrus.Infoln("starting grpc server")
	go func() {
		if err = srv.Serve(listen); err != nil {
			errC <- err
		}
	}()

	if err = <-errC; err != nil {
		logger.Log(logrus.DebugLevel, "could not serve", err)
	}

	logrus.Infoln("clean shutdown")
}

func main() {
	logger := logrus.New()

	srv := grpc.NewServer()
	dataMode := os.Getenv("STORAGE_MODE")
	shortener.RegisterUrlShortenerServer(srv, server.InitServer(dataMode, logger))
	execute(logger, srv)
}
