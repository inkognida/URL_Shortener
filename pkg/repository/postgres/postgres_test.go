package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
	"testing"
)

func testInit() *Data {
	logger := logrus.New()
	d := &Data{
		logger: logger,
	}

	dbUser := "hardella"
	dbPassword := "123"
	dbName := "postgres"

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		"localhost", "5432", dbUser, dbPassword, dbName)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		logger.Fatalln(err)

	}
	d.pool, err = pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		logger.Fatalln(err)

	}
	_, err = d.pool.Exec(context.Background(), `
   CREATE TABLE IF NOT EXISTS shorties (
      id BIGINT PRIMARY KEY,
      url TEXT NOT NULL,
      short_url TEXT NOT NULL UNIQUE
   )
`)
	if err != nil {
		logger.Fatalln(err)
	}
	return d
}

func TestData_GetUrl(t *testing.T) {
	data := testInit()

	_, err := data.pool.Exec(context.Background(), insertShortenUrlByID,
		-1, "google.com", "0000000000")
	if err != nil {
		t.Log(err)
	}

	mustGet := "google.com"
	link, err := data.GetUrl(context.Background(), "0000000000")
	if err != nil {
		t.Errorf("failed: no such link %s", err.Error())
	}
	if mustGet != link {
		t.Errorf("failed: wrong link %s, must be %s", link, mustGet)
	}
}

func TestData_SaveUrl(t *testing.T) {
	data := testInit()

	link := "google.com"
	short, err := data.SaveUrl(context.Background(), link)
	if err != nil {
		t.Errorf("failed: could not save the link %s", err.Error())
	}

	mustGet, err := data.GetUrl(context.Background(), short)
	if err != nil {
		t.Errorf("failed: saved wrong %s %s", err.Error(), link)
	}
	if mustGet != link {
		t.Errorf("failed: saved wrong link %s, must be %s", link, mustGet)
	}
}