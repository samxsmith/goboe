package goboe

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func ReadFile(notePath string) string {
	b, err := ioutil.ReadFile(notePath)
	if err != nil {
		panic(fmt.Sprintf("cant read file %s: %s", notePath, err))
	}
	return string(b)
}

var (
	ErrFileNotFound = errors.New("404")
)

func PathToAbs(path string) (string, error) {
	if strings.HasPrefix(path, "~") {
		home, ok := os.LookupEnv("HOME")
		if ok {
			path = strings.Replace(path, "~", home, 1)
		}
	}

	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return path, nil
}
