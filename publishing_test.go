package main

import (
	"errors"
	"os"
	"testing"

	. "github.com/ONSdigital/dp-train-tests/assertions"
	"github.com/ONSdigital/dp-train-tests/collections"
	"github.com/ONSdigital/dp-train-tests/train"
	"github.com/ONSdigital/dp-train-tests/website"
	"github.com/ONSdigital/log.go/log"
	. "github.com/smartystreets/goconvey/convey"
)

var collectionName = "otttest"

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

func Test_beingTransaction(t *testing.T) {
	Convey("Given an instance of the Train is running", t, func() {
		theTrain := train.NewClient()

		Convey("When a begin new transaction request is sent", func() {
			transaction, err := theTrain.Begin()
			defer transaction.CleanUp()

			So(err, ShouldBeNil)
			So(transaction, ShouldNotBeNil)

			Convey("Then a new transaction is created", func() {
				So(transaction, ShouldCreateTransaction)
			})
		})
	})
}

func Test_healthCheck(t *testing.T) {
	Convey("Given a healthy instance of the Train is running", t, func() {
		theTrain := train.NewClient()

		Convey("When a health check request is sent", func() {
			status, err := theTrain.HealthCheck()

			Convey("Then a healthy status response is returned", func() {
				So(err, ShouldBeNil)
				So(status, ShouldNotBeNil)
				So(status.Message, ShouldEqual, "OK")
			})
		})
	})
}

func Test_SendManifest(t *testing.T) {
	Convey("Given a publishing transaction has been created", t, func() {
		theTrain := train.NewClient()

		tx := beginTransaction(theTrain)
		defer tx.CleanUp()

		Convey("When a valid send manifest request is sent", func() {

			colDir := collections.GetPath(collectionName)
			err := theTrain.SendManifest(tx.ID, colDir)
			Convey("Then the request is successful", func() {
				So(err, ShouldBeNil)
			})

			Convey("And the published files to copy are copied into the transaction content dir", func() {

				manifest, err := collections.GetManifest(collectionName)
				So(err, ShouldBeNil)

				for _, item := range manifest.FilesToCopy {
					src := website.GetContentPath(item.Source)
					target := tx.GetContentURI(item.Target)

					So(target, ShouldCopyManifestEntryToContentDir, src)
				}
			})
		})
	})
}

func beginTransaction(theTrain *train.Client) *train.Transaction {
	tx, err := theTrain.Begin()

	So(err, ShouldBeNil)
	So(tx, ShouldCreateTransaction)

	return tx
}
