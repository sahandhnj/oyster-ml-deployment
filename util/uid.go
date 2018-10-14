package util

import (
	"strings"

	"github.com/google/uuid"
)

func UUID() string {
	return uuid.New().String()
}

func MinUUID(uuid string) string {
	parts := strings.SplitN(uuid, "-", 2)
	if len(parts) > 0 {
		return parts[0]
	}

	return ""
}
