package main

import (
	"github.com/sirupsen/logrus"
	"os"
	"shortener/internal/grpc_api/server"
	"shortener/internal/http_api/service"
)

func main() {
	logger := logrus.New()

	storageMode := os.Getenv("STORAGE_MODE")
	serviceMode := os.Getenv("SERVICE_MODE")
	if serviceMode == "grpc" {
		server.Execute(logger, storageMode)
	} else { //http
		httpService.Execute(logger, storageMode)
	}
}
