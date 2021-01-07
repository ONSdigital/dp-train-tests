package collections

import (
	"io/ioutil"
	"path/filepath"
)

type Content struct {
	File       string
	URI        string
	PublishURI string
	Zipped     bool
}

func GetContent(reviewedDir, file string) (*Content, error) {
	uri, err := filepath.Rel(reviewedDir, file)
	if err != nil {
		return nil, err
	}

	uri = filepath.Join("/", uri)
	publishURI := uri

	isZip := filepath.Ext(file) == ".zip"
	if isZip && filepath.Base(uri) == "timeseries-to-publish.zip" {
		publishURI = filepath.Join(filepath.Dir(uri), "timeseries")
	}

	return &Content{
		File:       file,
		URI:        uri,
		PublishURI: publishURI,
		Zipped:     isZip,
	}, nil
}

func (c Content) GetFile() string {
	return c.File
}

func (c Content) GetURI() string {
	return c.URI
}

func (c Content) GetPublishURI() string {
	return c.PublishURI
}

func (c Content) IsZip() bool {
	return c.Zipped
}

func (c Content) GetData() ([]byte, error) {
	return ioutil.ReadFile(c.File)
}
