package main

import (
	"context"
	"os"

	"github.com/ONSdigital/log.go/log"
)

func main() {
	if err := run(); err != nil {
		log.Event(context.Background(), "unexpected error", log.Error(err), log.ERROR)
		os.Exit(1)
	}
}

func run() error {
	return nil
}
