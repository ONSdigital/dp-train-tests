package train

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/ONSdigital/dp-train-tests/collections"
)

type PublishContent interface {
	GetFile() string
	GetURI() string
	GetPublishURI() string
	IsZip() bool
}

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
func (c *Client) SendManifest(txID string, manifest collections.Manifest) error {
	b, err := json.Marshal(manifest)
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

func (c *Client) AddContent(txID string, content PublishContent) error {
	body, contentType, err := newMultipartUpload(content.GetFile())
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/publish?transactionId=%s&uri=%s&zip=%t", c.Host, txID, content.GetPublishURI(), content.IsZip()), body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", contentType)

	var result ResponseEntity
	if err := c.doReq(req, http.StatusOK, &result); err != nil {
		return err
	}

	return nil
}

func newMultipartUpload(file string) (io.Reader, string, error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	filename := filepath.Base(file)
	if filename == "timeseries-to-publish.zip" {

	}

	part, err := writer.CreateFormFile("file", filepath.Base(file))
	if err != nil {
		return nil, "", err
	}

	f, err := os.Open(file)
	if err != nil {
		return nil, "", err
	}

	defer f.Close()

	_, err = io.Copy(part, f)
	if err != nil {
		return nil, "", err
	}

	writer.Close()
	return body, writer.FormDataContentType(), nil
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
