package assertions

import (
	"crypto/md5"
	"io/ioutil"
)

func ShouldCopyManifestEntryToContentDir(actual interface{}, expected ...interface{}) string {
	target, ok := actual.(string)
	if !ok {
		return "expected string for actual arg"
	}

	if len(expected) == 0 {
		return "expected string for expected[0] arg"
	}

	src, ok := expected[0].(string)
	if !ok {
		return "expected source for actual arg"
	}

	targetBytes, err := ioutil.ReadFile(target)
	if err != nil {
		return "error reading transaction content json: " + err.Error()
	}

	srcBytes, err := ioutil.ReadFile(src)
	if err != nil {
		return "error reading website content json: " + err.Error()
	}

	targetMD5 := md5.Sum(targetBytes)
	srcMD5 := md5.Sum(srcBytes)

	if srcMD5 != targetMD5 {
		return "incorrect md5 hash for transaction content: " + target
	}

	return ""
}
