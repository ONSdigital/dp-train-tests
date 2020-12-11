package assertions

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ONSdigital/dp-train-tests/train"
)

func ShouldCreateTransaction(actual interface{}, expected ...interface{}) string {
	if result := shouldCreateTransactionDir(actual, expected...); result != "" {
		return result
	}

	if result := shouldCreateBackupDir(actual, expected...); result != "" {
		return result
	}

	if result := shouldCreateContentDir(actual, expected...); result != "" {
		return result
	}

	if result := shouldCreateTransactionJSON(actual, expected...); result != "" {
		return result
	}

	return ""
}

// shouldCreateTransactionDir is an custom Smarty Streets Convey assertion function verifying the expected publishing
// transaction directory has been created. A *train.Transaction is required for the "actual" parameter.
func shouldCreateTransactionDir(actual interface{}, expected ...interface{}) string {
	tx, ok := actual.(*train.Transaction)
	if !ok {
		return "expected *train.Transaction for actual arg"
	}

	if !dirExists(tx.GetPath()) {
		return "expected transaction directory but none found"
	}

	return ""
}

// shouldCreateBackupDir is an custom Smarty Streets Convey assertion function verifying the expected publishing
// transaction backup directory has been created. A *train.Transaction is required for the "actual" parameter.
func shouldCreateBackupDir(actual interface{}, expected ...interface{}) string {
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

// shouldCreateContentDir is an custom Smarty Streets Convey assertion function verifying the expected publishing
// transaction content directory has been created. A *train.Transaction is required for the "actual" parameter.
func shouldCreateContentDir(actual interface{}, expected ...interface{}) string {
	tx, ok := actual.(*train.Transaction)
	if !ok {
		return "expected *train.Transaction for actual arg"
	}

	if !dirExists(tx.GetContentDirPath()) {
		return "expected transaction content directory but not found"
	}

	return ""
}

// shouldCreateTransactionJSON is an custom Smarty Streets Convey assertion function verifying the expected publishing
// transaction json file has been created. A *train.Transaction is required for the "actual" parameter.
func shouldCreateTransactionJSON(actual interface{}, expected ...interface{}) string {
	tx, ok := actual.(*train.Transaction)
	if !ok {
		return "expected *train.Transaction for expected arg[0]"
	}

	path := filepath.Join(tx.GetPath(), "transaction.json")
	f, err := os.Open(path)
	if err != nil {
		return "unexpected error while reading transaction.json file " + err.Error()
	}

	defer f.Close()

	var txJson train.Transaction
	if err := json.NewDecoder(f).Decode(&txJson); err != nil {
		return "unexpected error while marshalling transaction.json file " + err.Error()
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
