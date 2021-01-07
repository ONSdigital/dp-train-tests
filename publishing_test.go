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

var testCollection = "otttest"

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
			tx, err := theTrain.Begin()
			defer tx.CleanUp()

			So(err, ShouldBeNil)
			So(tx, ShouldNotBeNil)

			Convey("Then a new transaction is created", func() {
				So(tx, ShouldCreateTransaction)
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
			col := collections.Get(testCollection)

			manifest, err := col.GetManifest()
			So(err, ShouldBeNil)

			err = theTrain.SendManifest(tx.ID, manifest)

			Convey("Then the request is successful", func() {
				So(err, ShouldBeNil)
			})

			Convey("And the files to copy are correctly copied into the transaction content dir", func() {
				for _, item := range manifest.FilesToCopy {
					src := website.GetContentPath(item.Source)
					target := tx.GetContentFilePath(item.Target)

					So(target, ShouldCopyManifestEntryToContentDir, src)
				}
			})
		})
	})
}

func Test_AddContent(t *testing.T) {
	Convey("Given a transaction has been created", t, func() {
		theTrain := train.NewClient()

		tx := beginTransaction(theTrain)
		//defer tx.CleanUp()

		c := collections.Get(testCollection)
		sendManifest(tx, c, theTrain)

		Convey("When the collection content has been added to the transaction", func() {
			addContentToTransaction(tx, c, theTrain)

			collectionContent, err := c.GetAllContent()
			So(err, ShouldBeNil)

			Convey("Then the transaction contains the expected number of files", func() {
				transactionContent, err := tx.GetAllContent()
				So(err, ShouldBeNil)
				So(len(transactionContent), ShouldEqual, len(collectionContent))
			})

			Convey("And each collection content file exists in the transaction content", func() {
				for _, file := range collectionContent {
					So(file, ContentFileExistsInTransaction, tx)
				}
			})

			Convey("And each file contains the expected content", func() {
				for _, colContentItem := range collectionContent {
					So(colContentItem, ContentIsEqual, tx.GetContent(colContentItem.URI))
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

func sendManifest(tx *train.Transaction, c collections.Collection, theTrain *train.Client) {
	manifest, err := c.GetManifest()
	So(err, ShouldBeNil)

	err = theTrain.SendManifest(tx.ID, manifest)
	So(err, ShouldBeNil)
}

func addContentToTransaction(tx *train.Transaction, c collections.Collection, theTrain *train.Client) {
	content, err := c.ContentToPublish()
	So(err, ShouldBeNil)

	for _, item := range content {
		err := theTrain.AddContent(tx.ID, item)
		So(err, ShouldBeNil)
	}
}
