FROM golang:latest
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go build -o main ./cmd/main.go
EXPOSE 8000 8000
ENTRYPOINT ["/app/main"]