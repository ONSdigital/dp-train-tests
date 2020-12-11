package collections

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ONSdigital/dp-train-tests/train"
)

func collectionsDir() string {
	return filepath.Join(os.ExpandEnv("$zebedee_root"), "zebedee", "collections")
}

func GetPath(collectionName string) string {
	return filepath.Join(collectionsDir(), collectionName)
}

func GetManifest(collectionName string) (*train.Manifest, error) {
	b, err := ioutil.ReadFile(filepath.Join(GetPath(collectionName), "manifest.json"))
	if err != nil {
		return nil, err
	}

	var m train.Manifest
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}

	return &m, nil
}
