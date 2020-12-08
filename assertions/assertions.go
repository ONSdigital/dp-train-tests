package assertions

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ONSdigital/dp-train-tests/train"
)

// ShouldCreateTransactionDir is an custom Smarty Streets Convey assertion function verifying the expected publishing
// transaction directory has been created. A *train.Transaction is required for the "actual" parameter.
func ShouldCreateTransactionDir(actual interface{}, expected ...interface{}) string {
	tx, ok := actual.(*train.Transaction)
	if !ok {
		return "expected *train.Transaction for actual arg"
	}

	if !dirExists(tx.GetPath()) {
		return "expected transaction directory but none found"
	}

	return ""
}

// ShouldCreateBackupDir is an custom Smarty Streets Convey assertion function verifying the expected publishing
// transaction backup directory has been created. A *train.Transaction is required for the "actual" parameter.
func ShouldCreateBackupDir(actual interface{}, expected ...interface{}) string {
	tx, ok := actual.(*train.Transaction)
	if !ok {
		return "expected *train.Transaction for actual arg"
	}

	path := filepath.Join(tx.GetPath(), "backup")
	if !dirExists(path) {
		return "expected transaction backup directory but not found"
	}

	return ""
}

// ShouldCreateContentDir is an custom Smarty Streets Convey assertion function verifying the expected publishing
// transaction content directory has been created. A *train.Transaction is required for the "actual" parameter.
func ShouldCreateContentDir(actual interface{}, expected ...interface{}) string {
	tx, ok := actual.(*train.Transaction)
	if !ok {
		return "expected *train.Transaction for actual arg"
	}

	path := filepath.Join(tx.GetPath(), "content")
	if !dirExists(path) {
		return "expected transaction content directory but not found"
	}

	return ""
}

// ShouldCreateTransactionJSON is an custom Smarty Streets Convey assertion function verifying the expected publishing
// transaction json file has been created. A *train.Transaction is required for the "actual" parameter.
func ShouldCreateTransactionJSON(actual interface{}, expected ...interface{}) string {
	tx, ok := actual.(*train.Transaction)
	if !ok {
		return "expected *train.Transaction for expected arg[0]"
	}

	path := filepath.Join(tx.GetPath(), "transaction.json")
	f, err := os.Open(path)
	if err != nil {
		return "unexpected error while reading transaction.json file"
	}

	defer f.Close()

	var txJson train.Transaction
	if err := json.NewDecoder(f).Decode(&txJson); err != nil {
		return "unexpected error while marshalling transaction.json file"
	}

	if txJson.ID != tx.ID {
		return fmt.Sprintf("incorrect transaction.ID expected %q but was %q", tx.ID, txJson.ID)
	}

	if txJson.Status != "started" {
		return fmt.Sprintf("incorrect transaction.status expected %q but was %q", "started", txJson.Status)
	}

	return ""
}

func dirExists(path string) bool {
	var info os.FileInfo
	var err error

	if info, err = os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return info.IsDir()
}
