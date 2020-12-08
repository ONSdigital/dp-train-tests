package train

import (
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

func (t *Transaction) GetPath() string {
	return filepath.Join(os.Getenv("zebedee_root"), "zebedee", "transactions", t.ID)
}

func (t *Transaction) CleanUp() error {
	if err := os.RemoveAll(t.GetPath()); err != nil {
		log.Event(nil, "removing transaction directory", log.Error(err), log.ERROR, log.Data{"dir": t.GetPath()})
		return err
	}

	return nil
}
