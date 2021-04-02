package utils

import (
	"errors"
	"fmt"
	"io/ioutil"
)

func ReadNote(notePath string) string {
	b, err := ioutil.ReadFile(notePath)
	if err != nil {
		panic(fmt.Sprintf("cant read file %s: %s", notePath, err))
	}
	return string(b)
}

var (
	ErrFileNotFound = errors.New("404")
)
