package data

import (
	"context"
	"github.com/sirupsen/logrus"
	"testing"
)

func TestData_GetUrl(t *testing.T) {
	data := Data{}
	logger := logrus.New()

	err := data.Init(context.Background(), logger)
	if err != nil {
		t.Log(err)
	}

	data.urls["UQGhNVlVeL"] = "google.com"
	mustGet := "google.com"
	link, err := data.GetUrl(context.Background(), "UQGhNVlVeL")
	if err != nil {
		t.Errorf("failed: no such link %s", err.Error())
	}
	if mustGet != link {
		t.Errorf("failed: wrong link %s, must be %s", link, mustGet)
	}
}

func TestData_Init(t *testing.T) {
	data := Data{}
	logger := logrus.New()

	err := data.Init(context.Background(), logger)
	if err != nil {
		t.Errorf("failed: could not init %s", err.Error())
	}
}

func TestData_SaveUrl(t *testing.T) {
	data := Data{}
	logger := logrus.New()

	err := data.Init(context.Background(), logger)
	if err != nil {
		t.Log(err)
	}

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