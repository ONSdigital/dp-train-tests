package main

import (
	"errors"
	"os"
	"testing"

	"github.com/ONSdigital/dp-train-tests/train"
	"github.com/ONSdigital/log.go/log"
)

func TestMain(m *testing.M) {
	log.Namespace = "dp-train-tests"

	if err := setUp(); err != nil {
		log.Event(nil, "error setting up", log.Error(err), log.ERROR)
		os.Exit(1)
	}

	m.Run()
}

func setUp() error {
	trainCli := train.NewClient()
	if _, err := trainCli.HealthCheck(); err != nil {
		return errors.New("train instance not running - please start your instance and try again")
	}
	return nil
}
