package util

import (
	"strings"

	"github.com/google/uuid"
)

func UUID() string {
	return uuid.New().String()
}

func MinUUID(uuid string) string {
	return strings.TrimSuffix(uuid, "-")
}
