package main

import (
	"testing"

	. "github.com/ONSdigital/dp-train-tests/assertions"
	"github.com/ONSdigital/dp-train-tests/train"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_beingTransaction(t *testing.T) {
	Convey("Given an instance of the Train is running", t, func() {
		theTrain := train.NewClient()

		Convey("When a begin new transaction request is sent", func() {
			transaction, err := theTrain.Begin()
			So(err, ShouldBeNil)
			So(transaction, ShouldNotBeNil)

			defer transaction.CleanUp()

			Convey("Then a new transaction is created", func() {
				So(transaction, ShouldCreateTransactionDir)
				So(transaction, ShouldCreateBackupDir)
				So(transaction, ShouldCreateContentDir)
				So(transaction, ShouldCreateTransactionJSON)
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
