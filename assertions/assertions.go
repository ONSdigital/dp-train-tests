package assertions

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ONSdigital/dp-train-tests/train"
)

func ShouldCreateTransaction(actual interface{}, expected ...interface{}) string {
	transaction, ok := actual.(*train.Transaction)
	if !ok {
		return "expected transaction"
	}

	transDir := filepath.Join(train.GetInstance().TransactionsDir, transaction.ID)
	fmt.Printf("transaction dir %q\n", transDir)

	if !dirExists(transDir) {
		return "expected transaction directory but none found"
	}

	backupDir := filepath.Join(transDir, "backup")
	fmt.Printf("backup dir %q\n", backupDir)

	if !dirExists(backupDir) {
		return "expected backup directory but none found"
	}

	contentDir := filepath.Join(transDir, "content")
	if !dirExists(contentDir) {
		return "expected content directory but none found"
	}

	transactionJson := filepath.Join(transDir, "transaction.json")
	f, err := os.Open(transactionJson)
	if err != nil {
		return "unexpected error while reading transaction.json file"
	}

	defer f.Close()

	var t train.Transaction
	if err := json.NewDecoder(f).Decode(&t); err != nil {
		return "unexpected error while marshalling transaction.json file"
	}

	if t.ID != transaction.ID {

	}

	return ""
}

func ShouldCreateTransactionDir(actual interface{}, expected ...interface{}) string {
	trainInstance, ok := actual.(*train.Instance)
	if !ok {
		return "expected *train.Instance for actual arg"
	}

	transaction, ok := expected[0].(*train.Transaction)
	if !ok {
		return "expected *train.Transaction for expected arg[0]"
	}

	dir := filepath.Join(trainInstance.TransactionsDir, transaction.ID)
	if !dirExists(dir) {
		return "expected transaction directory but none found"
	}

	return ""
}

func ShouldCreateBackupDir(actual interface{}, expected ...interface{}) string {
	trainInstance, ok := actual.(*train.Instance)
	if !ok {
		return "expected *train.Instance for actual arg"
	}

	transaction, ok := expected[0].(*train.Transaction)
	if !ok {
		return "expected *train.Transaction for expected arg[0]"
	}

	dir := filepath.Join(trainInstance.TransactionsDir, transaction.ID, "backup")
	if !dirExists(dir) {
		return "expected backup directory but none found"
	}

	return ""
}

func ShouldCreateContentDir(actual interface{}, expected ...interface{}) string {
	trainInstance, ok := actual.(*train.Instance)
	if !ok {
		return "expected *train.Instance for actual arg"
	}

	transaction, ok := expected[0].(*train.Transaction)
	if !ok {
		return "expected *train.Transaction for expected arg[0]"
	}

	dir := filepath.Join(trainInstance.TransactionsDir, transaction.ID, "content")
	if !dirExists(dir) {
		return "expected transaction content directory but none found"
	}

	return ""
}

func ShouldCreateTransactionJSON(actual interface{}, expected ...interface{}) string {
	trainInstance, ok := actual.(*train.Instance)
	if !ok {
		return "expected *train.Instance for actual arg"
	}

	tx, ok := expected[0].(*train.Transaction)
	if !ok {
		return "expected *train.Transaction for expected arg[0]"
	}

	path := filepath.Join(trainInstance.TransactionsDir, tx.ID, "transaction.json")
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
