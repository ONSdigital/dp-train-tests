package website

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func getPath() string {
	return filepath.Join(os.ExpandEnv("$zebedee_root"), "zebedee", "master")
}

func GetContentPath(uri string) string {
	return filepath.Join(getPath(), uri)
}

func GetContent(uri string) ([]byte, error) {
	return ioutil.ReadFile(filepath.Join(getPath(), uri))
}
