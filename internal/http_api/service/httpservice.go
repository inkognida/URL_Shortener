package httpService

import (
	"context"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"shortener/internal/http_api/handler"
	"syscall"
	"time"
)

// process функция для обработки состояния сервера
func process(logger *logrus.Logger, srv *http.Server) (<-chan error, error) {
	errC := make(chan error, 1)
	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-ctx.Done()
		logger.Println("Grateful shutdown started")
		ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer func() {
			stop()
			cancel()
			close(errC)
		}()
		srv.SetKeepAlivesEnabled(false)
		if err := srv.Shutdown(ctxTimeout); err != nil {
			errC <- err
		}
		logger.Println("Shutdown stopped")
	}()
	go func() {
		logger.Println("Listen and serve", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errC <- err
		}
	}()

	return errC, nil
}

// Execute функция для создания нового обработчка и сервера
func Execute(logger *logrus.Logger, dataMode string) {
	srv, err := handler.NewHandler(dataMode, logger)
	if err != nil {
		logger.Fatalln(err)
	}

	errC, err := process(logger, &http.Server{
		Addr:              ":8000",
		Handler:           srv,
	})

	if err != nil {
		logger.Fatalln("Can't run")
	}
	if err := <-errC; err != nil {
		logger.Fatalln("Error while execution", err)
	}
}
