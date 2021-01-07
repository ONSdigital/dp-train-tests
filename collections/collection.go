package collections

import "C"
import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	timeSeriesDir = "/timeseries/"
)

type Collection struct {
	Path         string
	ReviewedDir  string
	ManifestPath string
}

func RootDir() string {
	return filepath.Join(os.ExpandEnv("$zebedee_root"), "zebedee", "collections")
}

func Get(name string) Collection {
	return Collection{
		Path:         filepath.Join(RootDir(), name),
		ReviewedDir:  filepath.Join(RootDir(), name, "reviewed"),
		ManifestPath: filepath.Join(RootDir(), name, "manifest.json"),
	}
}

func (c Collection) GetManifest() (m Manifest, err error) {
	var b []byte
	b, err = ioutil.ReadFile(c.ManifestPath)
	if err != nil {
		return m, err
	}

	if err := json.Unmarshal(b, &m); err != nil {
		return m, err
	}

	return m, nil
}

func (c Collection) ContentToPublish() ([]*Content, error) {
	isVersionedPath, err := regexp.Compile("^.*\\/v\\d+\\/.+")
	if err != nil {
		return nil, err
	}

	results := make([]*Content, 0)

	err = filepath.Walk(c.ReviewedDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || strings.Contains(path, timeSeriesDir) || isVersionedPath.MatchString(path) {
			return nil
		}

		content, err := GetContent(c.ReviewedDir, path)
		if err != nil {
			return err
		}

		results = append(results, content)

		return nil
	})

	return results, err
}

func (c Collection) GetAllContent() ([]*Content, error) {
	results := make([]*Content, 0)

	err := filepath.Walk(c.ReviewedDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || info.Name() == ".DS_Store"{
			return nil
		}

		content, err := GetContent(c.ReviewedDir, path)
		if err != nil {
			return err
		}

		results = append(results, content)

		return nil
	})

	return results, err
}
