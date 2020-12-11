package train

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"time"
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

// Begin creates a new publishing transaction.
func (c *Client) Begin() (*Transaction, error) {
	r, err := http.NewRequest(http.MethodPost, c.Host+"/begin", nil)
	if err != nil {
		return nil, err
	}

	var result ResponseEntity
	if err := c.doReq(r, http.StatusOK, &result); err != nil {
		return nil, err
	}

	return result.Transaction, nil
}

// HealthCheck sends a healthcheck request and checks the status is successful.
func (c *Client) HealthCheck() (*HealthStatus, error) {
	req, err := http.NewRequest("GET", c.Host+"/health", nil)
	if err != nil {
		return nil, err
	}

	var status HealthStatus
	if err := c.doReq(req, http.StatusOK, &status); err != nil {
		return nil, err
	}

	return &status, nil
}

// SendManifest send a collection manifest
func (c *Client) SendManifest(txID, collectionDir string) error {
	b, err := ioutil.ReadFile(filepath.Join(collectionDir, "manifest.json"))
	if err != nil {
		return err
	}

	r, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/CommitManifest?transactionId=%s", c.Host, txID), bytes.NewBuffer(b))
	if err != nil {
		return nil
	}

	var result ResponseEntity
	if err := c.doReq(r, http.StatusOK, &result); err != nil {
		return err
	}

	return nil
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
