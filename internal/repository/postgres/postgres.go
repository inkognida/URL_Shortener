package postgres

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"math/rand"
	"os"
	"shortener/internal/pkg/utils"
	"time"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Data struct {
	pool *pgxpool.Pool
	logger *logrus.Logger
}

const (
	HOST = "database"
	PORT = "5432"
)

func (d *Data) Init(ctx context.Context, logger *logrus.Logger) error {
	d.logger = logger
	dbUser, exist := os.LookupEnv("POSTGRES_USER")
	if !exist {
		dbUser = "admin"
	}
	dbPassword, exist := os.LookupEnv("POSTGRES_PASSWORD")
	if !exist {
		dbPassword = "admin"
	}
	dbName, exist := os.LookupEnv("POSTGRES_DB")
	if !exist {
		dbName = "admin"
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		HOST, PORT, dbUser, dbPassword, dbName)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return errors.Errorf("%s %v", "failed to parse config", err)
	}

	d.pool, err = pgxpool.ConnectConfig(ctx, poolConfig)
	if err != nil {
		return errors.Errorf("%s %v", "failed to connect via config", err)
	}

	_, err = d.pool.Exec(context.Background(), `
   CREATE TABLE IF NOT EXISTS shorties (
      id BIGINT PRIMARY KEY,
      url TEXT NOT NULL,
      short_url TEXT NOT NULL UNIQUE
   )
`)
	if err != nil {
		return errors.Errorf("%s %v", "failed to create a table", err)
	}

	return nil
}

const selectOriginalUrl = `SELECT url FROM shorties WHERE short_url = $1`

func (d *Data) GetUrl(ctx context.Context, shortUrl string) (string, error) {
	rows, err := d.pool.Query(ctx, selectOriginalUrl, shortUrl)
	if err != nil {
		return "", errors.Errorf("%s %v", "pool.Query error while getting", err)
	}

	defer rows.Close()

	var originalUrl string
	if rows.Next() {
		err = rows.Scan(&originalUrl)
		if err != nil {
			return "", errors.Errorf("%s %v", "no such link", err)
		}
	}

	return originalUrl, nil
}

const selectExistShortenUrlByID = `SELECT id FROM shorties WHERE id = $1`

const insertShortenUrlByID = `INSERT INTO shorties (id, url, short_url) VALUES ($1, $2, $3)`

func (d *Data) SaveUrl(ctx context.Context, originalUrl string) (string, error) {
	rand.Seed(time.Now().UnixNano())
	for {
		id := rand.Int63()
		rows, err := d.pool.Query(ctx, selectExistShortenUrlByID, id)
		if err != nil {
			rows.Close()
			return "", errors.Errorf("%s %v", "pool.Query error while saving", err)
		}

		if !rows.Next() {
			shortUrl := utils.GenerateLink(id)

			_, err = d.pool.Exec(ctx, insertShortenUrlByID, id, originalUrl, shortUrl)
			if err != nil {
				rows.Close()
				return "", errors.Errorf("%s %v", "pool.Exec error while saving", err)
			}

			rows.Close()
			return shortUrl, nil
		}
	}
}
