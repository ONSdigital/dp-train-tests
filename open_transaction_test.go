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

		Convey("When Begin is invoked", func() {
			transaction, err := theTrain.Begin()
			So(err, ShouldBeNil)

			Convey("Then a new transaction is created", func() {
				t := train.GetInstance()

				So(transaction, ShouldNotBeNil)
				So(t, ShouldCreateTransactionDir, transaction)
				So(t, ShouldCreateBackupDir, transaction)
				So(t, ShouldCreateContentDir, transaction)
				So(t, ShouldCreateTransactionJSON, transaction)
			})
		})
	})
}
