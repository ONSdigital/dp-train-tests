package main

import (
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/ONSdigital/dp-train-tests/train"
	"github.com/ONSdigital/log.go/log"
)

func TestMain(m *testing.M) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	t, err := startUp()
	if err != nil {
		log.Event(nil, "error", log.Error(err), log.ERROR)
		os.Exit(1)
	}

	go func() {
		m.Run()
		train.GetInstance().CompletedC <- true
	}()

	select {
	case <-t.CompletedC:
		log.Event(nil, "completed tests shutting down", log.INFO)
		t.Stop()
	case err := <-t.ErrC:
		log.Event(nil, "train runner returned an error shutting down", log.Error(err), log.ERROR)
		t.Stop()
	case s := <-signals:
		log.Event(nil, "signal intercepted shutting down", log.INFO, log.Data{"signal": s.String()})
		t.Stop()
	}
}

func startUp() (*train.Instance, error) {
	t, err := train.NewRunner()
	if err != nil {
		return nil, err
	}

	if err := t.Start(); err != nil {
		return nil, err
	}

	return t, nil
}
