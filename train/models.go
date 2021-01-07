package train

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ONSdigital/log.go/log"
)

type ResponseEntity struct {
	Transaction *Transaction `json:"transaction"`
}

type Transaction struct {
	ID         string        `json:"id"`
	Status     string        `json:"status"`
	StartDate  string        `json:"startDate"`
	EndDate    string        `json:"startDate"`
	UriInfos   []interface{} `json:"uriInfos"`
	UriDeletes []interface{} `json:"uriDeletes"`
	Errors     []interface{} `json:"errors"`
	Files      interface{}   `json:"files"`
}

type HealthStatus struct {
	Message string `json:"message"`
}

type Content struct {
	File  string
	URI   string
	IsZip bool
}

func (t *Transaction) GetPath() string {
	return filepath.Join(os.Getenv("zebedee_root"), "zebedee", "transactions", t.ID)
}

func (t *Transaction) GetContentDirPath() string {
	return filepath.Join(t.GetPath(), "content")
}

func (t *Transaction) GetContentFilePath(uri string) string {
	return filepath.Join(t.GetContentDirPath(), uri)
}

func (t *Transaction) GetContent(uri string) *Content {
	return &Content{
		File:  filepath.Join(t.GetContentDirPath(), uri),
		URI:   uri,
		IsZip: filepath.Ext(uri) == ".zip",
	}
}

func (t *Transaction) ContentExists(uri string) bool {
	path := t.GetContentFilePath(uri)

	if _, err := os.Stat(path); err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	} else {
		fmt.Printf("error checking file exists: file: %s, error: %s\n", path, err.Error())
		return false
	}
}

func (t *Transaction) CleanUp() error {
	if err := os.RemoveAll(t.GetPath()); err != nil {
		log.Event(nil, "removing transaction directory", log.Error(err), log.ERROR, log.Data{"dir": t.GetPath()})
		return err
	}

	return nil
}

func (t *Transaction) GetAllContent() ([]os.FileInfo, error) {
	infos := make([]os.FileInfo, 0)
	err := filepath.Walk(t.GetContentDirPath(), func(path string, info os.FileInfo, err error) error {

		if info.IsDir() || info.Name() == ".DS_Store" {
			return nil
		}

		infos = append(infos, info)
		return nil
	})

	return infos, err
}

func (c *Content) GetData() ([]byte, error) {
	return ioutil.ReadFile(c.File)
}
