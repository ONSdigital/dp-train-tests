package train

import (
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

type FileCopy struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

type FileDelete struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

type Manifest struct {
	FilesToCopy  []*FileCopy   `json:"filesToCopy"`
	UrisToDelete []*FileDelete `json:"uriToDelete"`
}

func (t *Transaction) GetPath() string {
	return filepath.Join(os.Getenv("zebedee_root"), "zebedee", "transactions", t.ID)
}

func (t *Transaction) GetContentDirPath() string {
	return filepath.Join(t.GetPath(), "content")
}

func (t *Transaction) GetContentURI(uri string) string {
	return filepath.Join(t.GetContentDirPath(), uri)
}

func (t *Transaction) GetContent(uri string) ([]byte, error) {
	return ioutil.ReadFile( filepath.Join(t.GetContentDirPath(), uri))
}

func (t *Transaction) CleanUp() error {
	if err := os.RemoveAll(t.GetPath()); err != nil {
		log.Event(nil, "removing transaction directory", log.Error(err), log.ERROR, log.Data{"dir": t.GetPath()})
		return err
	}

	return nil
}
