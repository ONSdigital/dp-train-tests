package train

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Client struct {
	Host    string
	HttpCli http.Client
}

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

func NewClient() *Client {
	return &Client{
		Host:    "http://localhost:8084",
		HttpCli: http.Client{Timeout: time.Second * 5},
	}
}

func (c *Client) Begin() (*Transaction, error) {
	r, err := http.NewRequest(http.MethodPost, c.Host+"/begin", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.HttpCli.Do(r)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("begin transaction: incorrect response status expected 200 but was %d", resp.StatusCode)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var entity ResponseEntity
	if err := json.Unmarshal(b, &entity); err != nil {
		return nil, err
	}

	return entity.Transaction, nil
}
