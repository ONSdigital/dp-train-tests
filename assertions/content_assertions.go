package assertions

import (
	"crypto/md5"
	"fmt"

	"github.com/ONSdigital/dp-train-tests/collections"
	"github.com/ONSdigital/dp-train-tests/train"
)

func ContentFileExistsInTransaction(actual interface{}, expected ...interface{}) string {
	content, ok := actual.(*collections.Content)
	if !ok {
		return "expected *collections.Content for actual arg"
	}

	tx, ok := expected[0].(*train.Transaction)
	if !ok {
		return "expected *train.Transaction for expected arg"
	}

	if !tx.ContentExists(content.URI) {
		return fmt.Sprintf("expected file in transaction content but not found, uri: %s", content.URI)
	}

	return ""
}

func ContentIsEqual(actual interface{}, expected ...interface{}) string {
	collectionContent, ok := actual.(*collections.Content)
	if !ok {
		return "expected *collections.Content for actual arg"
	}

	transactionContent, ok := expected[0].(*train.Content)
	if !ok {
		return "expected *train.Content for expected arg 0"
	}

	expectedData, err := collectionContent.GetData()
	if err != nil {
		return fmt.Sprintf("error getting collection content bytes, uri: %s error: %s", collectionContent.URI, err.Error())
	}

	actualData, err := transactionContent.GetData()
	if err != nil {
		return fmt.Sprintf("error getting transaction content bytes, uri: %s error: %s", transactionContent.URI, err.Error())
	}

	if md5.Sum(expectedData) != md5.Sum(actualData) {
		return "incorrect md5 hash for transaction content: " + transactionContent.URI
	}

	return ""
}