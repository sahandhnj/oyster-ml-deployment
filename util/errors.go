package util

import "errors"

const (
	ErrNotFound = "Object was not found"
)

func GetError(errorMessage string) error {
	return errors.New(errorMessage)
}
