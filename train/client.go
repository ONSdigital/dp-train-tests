package train

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/ONSdigital/log.go/log"
)

type Client struct {
	Host    string
	HttpCli http.Client
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

	var result ResponseEntity
	if err := c.doReq(r, http.StatusOK, &result); err != nil {
		return nil, err
	}

	log.Event(nil, "begin transaction request completed successfully", log.INFO, log.Data{"id": result.Transaction.ID})
	return result.Transaction, nil
}

func (c *Client) HealthCheck() (*HealthStatus, error) {
	req, err := http.NewRequest("GET", "http://localhost:8084/health", nil)
	if err != nil {
		return nil, err
	}

	var status HealthStatus
	if err := c.doReq(req, http.StatusOK, &status); err != nil {
		return nil, err
	}

	log.Event(nil, "The-Train HealthCheck request successful", log.INFO)
	return &status, nil
}

func (c *Client) doReq(req *http.Request, expectedStatus int, entity interface{}) error {
	resp, err := c.HttpCli.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != expectedStatus {
		return fmt.Errorf("incorrect response status expected 200 but was %d", resp.StatusCode)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(b, entity); err != nil {
		return err
	}

	return nil
}


